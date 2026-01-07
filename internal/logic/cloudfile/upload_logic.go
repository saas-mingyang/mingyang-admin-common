package cloudfile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/tenantctx"
	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"github.com/saas-mingyang/mingyang-admin-common/utils/sonyflake"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang-admin-simple-admin-file/ent"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
	"mingyang-admin-simple-admin-file/internal/utils/cloud"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"
	"mingyang-admin-simple-admin-file/internal/utils/filex"
)

// 定义分片大小
const (
	smallFileThreshold = 10 * 1024 * 1024 // 10MB
	chunkSize          = 10 * 1024 * 1024 // 10MB
)

// ===================== 进度管理器 =====================

// UploadProgress 上传进度结构体
// UploadProgress 文件上传进度信息结构体
// 用于实时跟踪和记录文件上传的进度状态
type UploadProgress struct {
	// UploadID 上传任务唯一标识符
	// 使用雪花算法生成，用于查询和跟踪特定上传任务
	UploadID uint64 `json:"uploadId"`

	// FileName 原始文件名（用户上传的文件名）
	// 包含扩展名，如：example.jpg
	FileName string `json:"fileName"`

	// TotalSize 文件总大小（字节）
	// 完整文件的字节数，用于计算上传百分比
	TotalSize int64 `json:"totalSize"`

	// Uploaded 已上传字节数
	// 实时更新的已成功上传的字节数
	Uploaded int64 `json:"uploaded"`

	// TotalParts 总分片数（仅分片上传时有效）
	// 大文件分片上传时，文件被分割成的总块数
	// 计算方式：ceil(TotalSize / chunkSize)
	TotalParts int `json:"totalParts"`

	// CurrentPart 当前正在上传的分片序号
	// 从1开始计数，表示当前正在上传第几个分片
	CurrentPart int `json:"currentPart"`

	// Percentage 上传百分比（0.0 - 100.0）
	// 计算方式：(Uploaded / TotalSize) * 100
	// 前端可用于显示进度条
	Percentage float64 `json:"percentage"`

	// Speed 上传速度（KB/s，千字节每秒）
	// 实时计算的平均上传速度，用于预估剩余时间
	// 计算方式：Uploaded / 已用时间 / 1024
	Speed float64 `json:"speed"`

	// Status 上传状态
	// 枚举值：preparing(准备中)、uploading(上传中)、completed(已完成)、failed(失败)
	// 状态流转：preparing → uploading → completed/failed
	Status string `json:"status"`

	// StartTime 上传开始时间
	// 记录上传任务开始的时间戳，用于计算已用时间和上传速度
	StartTime time.Time `json:"startTime"`

	// UserID 用户唯一标识符
	// 上传文件的用户ID，用于按用户查询上传任务
	UserID string `json:"userId"`

	// Provider 云存储服务提供商
	// 枚举值：aws_s3、aliyun_oss、tencent_cos、qiniu_kodo等
	// 表示文件上传到哪个云存储服务
	Provider string `json:"provider"`

	// Bucket 云存储桶名称
	// 文件存储的目标存储桶，用于区分不同的存储空间
	Bucket string `json:"bucket"`

	// Key 文件在云存储中的唯一路径标识
	// 格式示例：2024-01-15/tenant_id/file_type/filename.ext
	// 包含目录结构和文件名，确保文件唯一性
	Key string `json:"key"`
}

// ProgressManager 进度管理器
type ProgressManager struct {
	sync.RWMutex
	progressMap map[uint64]*UploadProgress // key: uploadId -> progress
}

var progressManager = &ProgressManager{
	progressMap: make(map[uint64]*UploadProgress),
}

// AddProgress 添加进度记录
func (pm *ProgressManager) AddProgress(uploadId uint64, progress *UploadProgress) {
	pm.Lock()
	defer pm.Unlock()
	pm.progressMap[uploadId] = progress
	logx.Infow("添加进度记录", logx.Field("uploadId", uploadId), logx.Field("fileName", progress.FileName))
}

// UpdateProgress 更新进度
func (pm *ProgressManager) UpdateProgress(uploadId uint64, uploaded int64, currentPart int) {
	pm.Lock()
	defer pm.Unlock()

	if progress, exists := pm.progressMap[uploadId]; exists {
		progress.Uploaded = uploaded
		progress.CurrentPart = currentPart
		if progress.TotalSize > 0 {
			progress.Percentage = float64(uploaded) / float64(progress.TotalSize) * 100
		}

		// 计算上传速度
		elapsed := time.Since(progress.StartTime).Seconds()
		if elapsed > 0 {
			progress.Speed = float64(uploaded) / elapsed / 1024 // KB/s
		}
	}
}

// UpdateStatus 更新状态
func (pm *ProgressManager) UpdateStatus(uploadId uint64, status string) {
	pm.Lock()
	defer pm.Unlock()

	if progress, exists := pm.progressMap[uploadId]; exists {
		progress.Status = status
		if status == "completed" {
			progress.Uploaded = progress.TotalSize
			progress.Percentage = 100
		}
	}
}

// GetProgress 获取进度
func (pm *ProgressManager) GetProgress(uploadId uint64) (*UploadProgress, bool) {
	pm.RLock()
	defer pm.RUnlock()

	progress, exists := pm.progressMap[uploadId]
	return progress, exists
}

// DeleteProgress 删除进度记录
func (pm *ProgressManager) DeleteProgress(uploadId uint64) {
	pm.Lock()
	defer pm.Unlock()
	delete(pm.progressMap, uploadId)
}

// CleanupOldProgress 清理旧的进度记录（超过24小时）
func (pm *ProgressManager) CleanupOldProgress() {
	pm.Lock()
	defer pm.Unlock()

	now := time.Now()
	for uploadId, progress := range pm.progressMap {
		if now.Sub(progress.StartTime) > 24*time.Hour {
			delete(pm.progressMap, uploadId)
		}
	}
}

// ===================== 带进度的Reader =====================

// progressReader 带进度的Reader，实现io.ReadSeeker
type progressReader struct {
	reader     io.ReadSeeker
	totalRead  int64
	onProgress func(int64)
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.totalRead += int64(n)
		if r.onProgress != nil {
			r.onProgress(r.totalRead)
		}
	}
	return n, err
}

func (r *progressReader) Seek(offset int64, whence int) (int64, error) {
	// 先重置totalRead
	if whence == io.SeekStart && offset == 0 {
		r.totalRead = 0
	}

	// 调用底层的Seek
	newPos, err := r.reader.Seek(offset, whence)
	if err != nil {
		return newPos, err
	}

	// 更新totalRead
	if whence == io.SeekCurrent {
		r.totalRead += offset
	} else if whence == io.SeekStart {
		r.totalRead = offset
	}
	// 对于io.SeekEnd，我们无法确定当前读取位置，所以不更新totalRead

	return newPos, nil
}

// ===================== UploadLogic =====================

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadLogic(r *http.Request, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		r:      r,
	}
}

// getFileTypeByExtension 根据文件后缀判断文件类型
func getFileTypeByExtension(filename string) string {
	// 获取文件后缀（不包含点）
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "other"
	}

	// 去掉点
	ext = strings.TrimPrefix(ext, ".")

	// 图片类型
	imageExts := []string{
		"jpg", "jpeg", "png", "gif", "bmp", "webp",
		"tiff", "tif", "svg", "ico", "psd", "raw",
	}

	// 视频类型
	videoExts := []string{
		"mp4", "avi", "mov", "mkv", "flv", "wmv",
		"m4v", "mpg", "mpeg", "3gp", "webm", "ogg",
	}

	// 音频类型
	audioExts := []string{
		"mp3", "wav", "aac", "flac", "wma", "m4a",
		"ogg", "oga", "opus", "mid", "midi",
	}

	// APK类型
	if ext == "apk" {
		return "apk"
	}

	// 检查图片
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return "image"
		}
	}

	// 检查视频
	for _, vidExt := range videoExts {
		if ext == vidExt {
			return "video"
		}
	}

	// 检查音频
	for _, audExt := range audioExts {
		if ext == audExt {
			return "audio"
		}
	}

	return "other"
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if size < KB {
		return fmt.Sprintf("%d B", size)
	} else if size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	} else if size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
}

func (l *UploadLogic) Upload() (resp *types.CloudFileInfoResp, err error) {
	tenantId := tenantctx.GetTenantIDFromCtx(l.ctx)
	if _, ok := l.svcCtx.CloudStorage.Service[tenantId]; !ok {
		err = cloud.AddTenantCloudServiceGroup(l.svcCtx.DB, l.svcCtx.CloudStorage, tenantId)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errorx.NewCodeInternalError("storage_provider.StorageProviderNotExist")
			}
			return nil, errorx.NewCodeInternalError("storage_provider.failedLoadProviderConfig")
		}
	}

	err = l.r.ParseMultipartForm(l.svcCtx.Config.UploadConf.MaxVideoSize)
	if err != nil {
		logx.Error("fail to parse the multipart form")
		return nil, errorx.NewCodeInvalidArgumentError(
			"file.parseFormFailed")
	}

	file, handler, err := l.r.FormFile("file")
	if err != nil {
		logx.Error("the value of file cannot be found")
		return nil, errorx.NewCodeInvalidArgumentError("file.parseFormFailed")
	}
	defer file.Close()

	// judge if the suffix is legal
	dotIndex := strings.LastIndex(handler.Filename, ".")
	if dotIndex == -1 {
		logx.Errorw("reject the file which does not have suffix")
		return nil, errorx.NewCodeInvalidArgumentError("file.wrongTypeError")
	}

	fileName, fileSuffix := handler.Filename[:dotIndex], handler.Filename[dotIndex+1:]
	fileUUID := sonyflake.NextID()
	storeFileName := fmt.Sprint(fileUUID) + "." + fileSuffix
	userId := l.ctx.Value("userId").(string)

	// 使用后缀判断文件类型
	fileType := getFileTypeByExtension(handler.Filename)

	logx.Infow("文件上传信息",
		logx.Field("文件名", handler.Filename),
		logx.Field("文件大小", formatFileSize(handler.Size)),
		logx.Field("文件类型", fileType),
		logx.Field("Content-Type", handler.Header.Get("Content-Type")))

	// judge if the file size is over max size
	// 判断文件大小是否超过设定值
	err = filex.CheckOverSize(l.ctx, l.svcCtx, fileType, handler.Size)
	if err != nil {
		logx.Errorw("the file is over size", logx.Field("type", fileType),
			logx.Field("userId", userId), logx.Field("size", handler.Size),
			logx.Field("fileName", handler.Filename))
		return nil, err
	}

	var provider string
	if l.r.MultipartForm.Value["provider"] != nil && l.svcCtx.CloudStorage.Service[tenantId].CloudStorage[l.r.MultipartForm.Value["provider"][0]] != nil {
		provider = l.r.MultipartForm.Value["provider"][0]
	} else {
		provider = l.svcCtx.CloudStorage.Service[tenantId].DefaultProvider
	}

	var fileTagId uint64
	if l.r.MultipartForm.Value["tagId"] != nil && l.r.MultipartForm.Value["tagId"][0] != "" {
		tagId, err := strconv.Atoi(l.r.MultipartForm.Value["tagId"][0])
		if err != nil {
			return nil, errorx.NewCodeInvalidArgumentError("wrong tag ID")
		}

		fileTagId = uint64(tagId)
	}

	relativeSrc := fmt.Sprintf("%s/%s/%s/%s",
		datetime.FormatTimeToStr(time.Now(), "yyyy-mm-dd"),
		fmt.Sprint(tenantId),
		fileType,
		storeFileName)

	// 生成上传ID（uploadId）
	uploadId := fileUUID

	// 计算总块数
	totalParts := int((handler.Size + chunkSize - 1) / chunkSize) // 向上取整

	// 创建进度记录
	progress := &UploadProgress{
		UploadID:    uploadId,
		FileName:    handler.Filename,
		TotalSize:   handler.Size,
		Uploaded:    0,
		TotalParts:  totalParts,
		CurrentPart: 0,
		Percentage:  0,
		Speed:       0,
		Status:      "preparing",
		StartTime:   time.Now(),
		UserID:      userId,
		Provider:    provider,
		Bucket:      l.svcCtx.CloudStorage.Service[tenantId].ProviderData[provider].Bucket,
		Key:         relativeSrc,
	}

	// 添加进度记录
	progressManager.AddProgress(uploadId, progress)

	// 启动一个goroutine定期清理旧的进度记录
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			progressManager.CleanupOldProgress()
		}
	}()

	// 开始上传
	progressManager.UpdateStatus(uploadId, "uploading")

	// 根据文件大小选择上传方式
	var url string
	if handler.Size < smallFileThreshold {
		url, err = l.UploadToProviderSimple(file, relativeSrc, provider, tenantId, uploadId)
	} else {
		url, err = l.UploadToProviderMultipart(file, relativeSrc, provider, tenantId, uploadId, handler.Size)
	}

	// 延迟删除进度记录（5分钟后）
	defer func() {
		go func() {
			time.Sleep(5 * time.Minute)
			progressManager.DeleteProgress(uploadId)
		}()
	}()

	if err != nil {
		progressManager.UpdateStatus(uploadId, "failed")
		logx.Errorw("上传到云存储失败",
			logx.Field("error", err),
			logx.Field("fileName", handler.Filename),
			logx.Field("size", handler.Size),
			logx.Field("provider", provider))
		return nil, err
	}

	// 更新状态为完成
	progressManager.UpdateStatus(uploadId, "completed")

	logx.Infow("文件上传成功",
		logx.Field("url", url),
		logx.Field("fileName", handler.Filename),
		logx.Field("size", handler.Size),
		logx.Field("provider", provider),
		logx.Field("uploadId", uploadId))

	service := l.svcCtx.CloudStorage.Service[tenantId]
	storageProvider := service.ProviderData[provider]

	// store to database
	query := l.svcCtx.DB.CloudFile.Create().
		SetID(fileUUID).
		SetName(fileName).
		SetFileType(filex.ConvertFileTypeToUint8(fileType)).
		SetStorageProvidersID(storageProvider.Id).
		SetURL(url).
		SetSize(uint64(handler.Size)).
		SetUserID(userId)

	if fileTagId != 0 {
		query = query.AddTagIDs(fileTagId)
	}

	data, err := query.Save(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, nil)
	}

	logic := NewGetCloudFileDownloadUrlLogic(l.ctx, l.svcCtx)

	downloadUrl, err := logic.GetCloudFileDownloadUrl(&types.BaseIDInfo{Id: &data.ID})

	if err != nil {
		return nil, err
	}
	return &types.CloudFileInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  i18n.Success,
		},
		Data: types.CloudFileInfo{
			BaseIDInfo: types.BaseIDInfo{
				Id:        &data.ID,
				CreatedAt: pointy.GetPointer(data.CreatedAt.UnixMilli()),
			},
			State:       pointy.GetPointer(data.State),
			Name:        pointy.GetPointer(data.Name),
			Url:         pointy.GetPointer(*downloadUrl.Data.Url),
			RelativeSrc: pointy.GetPointer(relativeSrc),
			Size:        pointy.GetPointer(data.Size),
			FileType:    pointy.GetPointer(data.FileType),
			UserId:      pointy.GetPointer(data.UserID),
		},
	}, nil
}

// UploadToProviderSimple 普通上传（带超时和重试）
func (l *UploadLogic) UploadToProviderSimple(file multipart.File, fileName, provider string, tenantId uint64, uploadId uint64) (url string, err error) {
	logx.Infow("普通上传",
		logx.Field("文件名", fileName),
		logx.Field("提供商", provider),
		logx.Field("uploadId", uploadId))

	if client, ok := l.svcCtx.CloudStorage.Service[tenantId].CloudStorage[provider]; ok {
		// 设置较长的超时时间
		ctx, cancel := context.WithTimeout(l.ctx, 10*time.Minute)
		defer cancel()

		// 创建带进度的ReadSeeker
		// multipart.File已经实现了io.ReadSeeker，所以我们可以直接使用
		progressReader := &progressReader{
			reader: file,
			onProgress: func(readBytes int64) {
				progressManager.UpdateProgress(uploadId, readBytes, 1)
			},
		}

		_, err := client.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(l.svcCtx.CloudStorage.Service[tenantId].ProviderData[provider].Bucket),
			Key:    aws.String(fileName),
			Body:   progressReader,
		})
		if err != nil {
			fmt.Printf("failed to upload object %s\n", err)
			logx.Errorw("failed to upload object", logx.Field("detail", err))

			// 检查错误类型
			var aerr awserr.Error
			if errors.As(err, &aerr) {
				logx.Errorw("AWS错误详情",
					logx.Field("Code", aerr.Code()),
					logx.Field("Message", aerr.Message()))

				if aerr.Code() == request.CanceledErrorCode {
					// 检查是否超时
					if errors.Is(err, context.DeadlineExceeded) {
						return "", errorx.NewCodeInternalError("上传超时，请重试或使用分片上传")
					}
					return "", errorx.NewCodeInternalError("上传被取消")
				}
			}
			return "", errorx.NewCodeInternalError("上传失败")
		}
		logx.Infow("普通上传成功", logx.Field("fileName", fileName))
		return fileName, nil
	}
	return "", fmt.Errorf("云存储客户端未找到: %s", provider)
}

// UploadToProviderMultipart 分片上传实现
func (l *UploadLogic) UploadToProviderMultipart(file multipart.File, fileName, provider string, tenantId uint64, uploadId uint64, totalSize int64) (url string, err error) {
	logx.Infow("开始分片上传",
		logx.Field("文件名", fileName),
		logx.Field("提供商", provider),
		logx.Field("分片大小", chunkSize),
		logx.Field("uploadId", uploadId),
		logx.Field("总大小", formatFileSize(totalSize)))

	if client, ok := l.svcCtx.CloudStorage.Service[tenantId].CloudStorage[provider]; ok {
		bucket := l.svcCtx.CloudStorage.Service[tenantId].ProviderData[provider].Bucket

		// 1. 初始化分片上传
		createResp, err := client.CreateMultipartUploadWithContext(l.ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileName),
		})
		if err != nil {
			logx.Errorw("创建分片上传失败",
				logx.Field("error", err),
				logx.Field("bucket", bucket),
				logx.Field("key", fileName))
			return "", err
		}

		uploadS3Id := createResp.UploadId
		logx.Infow("分片上传已创建",
			logx.Field("UploadId", *uploadS3Id))

		// 2. 分片上传
		var completedParts []*s3.CompletedPart
		partNumber := int64(1)
		buffer := make([]byte, chunkSize)

		// 记录已上传字节数
		totalUploaded := int64(0)

		for {
			// 读取一个分片
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				logx.Errorw("读取文件分片失败",
					logx.Field("error", err),
					logx.Field("partNumber", partNumber))
				return "", err
			}

			// 如果读取到数据为0，说明文件已读完
			if n == 0 {
				break
			}

			// 更新进度
			totalUploaded += int64(n)
			progressManager.UpdateProgress(uploadId, totalUploaded, int(partNumber))

			// 记录当前进度日志
			if progress, exists := progressManager.GetProgress(uploadId); exists {
				logx.Infow("上传分片",
					logx.Field("partNumber", partNumber),
					logx.Field("size", n),
					logx.Field("uploaded", totalUploaded),
					logx.Field("progress", fmt.Sprintf("%.1f%%", progress.Percentage)),
					logx.Field("speed", fmt.Sprintf("%.1f KB/s", progress.Speed)))
			}

			partResp, err := client.UploadPartWithContext(l.ctx, &s3.UploadPartInput{
				Bucket:     aws.String(bucket),
				Key:        aws.String(fileName),
				UploadId:   uploadS3Id,
				PartNumber: aws.Int64(partNumber),
				Body:       bytes.NewReader(buffer[:n]),
			})

			if err != nil {
				logx.Errorw("上传分片失败",
					logx.Field("error", err),
					logx.Field("partNumber", partNumber))

				// 尝试中止上传
				_, abortErr := client.AbortMultipartUploadWithContext(l.ctx, &s3.AbortMultipartUploadInput{
					Bucket:   aws.String(bucket),
					Key:      aws.String(fileName),
					UploadId: uploadS3Id,
				})
				if abortErr != nil {
					logx.Errorw("中止分片上传失败",
						logx.Field("error", abortErr))
				}
				return "", err
			}

			// 保存完成的分片信息
			completedParts = append(completedParts, &s3.CompletedPart{
				ETag:       partResp.ETag,
				PartNumber: aws.Int64(partNumber),
			})

			partNumber++

			// 如果遇到EOF，说明文件已读完
			if err == io.EOF {
				break
			}
		}

		// 3. 完成分片上传
		completeResp, err := client.CompleteMultipartUploadWithContext(l.ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(fileName),
			UploadId: uploadS3Id,
			MultipartUpload: &s3.CompletedMultipartUpload{
				Parts: completedParts,
			},
		})

		if err != nil {
			logx.Errorw("完成分片上传失败",
				logx.Field("error", err))
			return "", err
		}

		logx.Infow("分片上传完成",
			logx.Field("Location", *completeResp.Location),
			logx.Field("总块数", len(completedParts)))

		return fileName, nil
	}

	return "", fmt.Errorf("云存储客户端未找到: %s", provider)
}
