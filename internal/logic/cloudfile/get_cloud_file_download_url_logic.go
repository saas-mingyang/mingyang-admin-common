package cloudfile

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang.com/admin-simple-admin-file/ent/cloudfile"
	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"
	"strings"
	"time"
)

type GetCloudFileDownloadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCloudFileDownloadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileDownloadUrlLogic {
	return &GetCloudFileDownloadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCloudFileDownloadUrlLogic) GetCloudFileDownloadUrl(req *types.BaseIDInfo) (resp *types.CloudFileInfoResp, err error) {
	file, err := l.svcCtx.DB.CloudFile.
		Query().
		Where(cloudfile.ID(*req.Id)).
		WithStorageProviders().
		First(l.ctx)
	if err != nil || file == nil {
		return nil, errorx.NewCodeInvalidArgumentError(err.Error())
	}

	if !file.IsDownloaded {
		return nil, errorx.NewCodeInvalidArgumentError("file.fileNotUploaded")
	}

	providers := file.Edges.StorageProviders
	if providers == nil {
		return nil, errorx.NewCodeInvalidArgumentError(err.Error())
	}
	privateURL, err := GetPrivateURLExact(providers.CdnURL, file.URL, providers.SecretID, providers.SecretKey)
	if err != nil {
		return nil, errorx.NewCodeInternalError(err.Error())
	}

	resp = &types.CloudFileInfoResp{}

	id := file.ID
	resp.Data = types.CloudFileInfo{
		Url:        &privateURL,
		Size:       &file.Size,
		Name:       &file.Name,
		FileType:   &file.FileType,
		UserId:     &file.UserID,
		ProviderId: &providers.ID,
		BaseIDInfo: types.BaseIDInfo{
			Id: &id,
		},
	}

	return resp, nil
}

// GetPrivateURLExact 生成私有文件访问URL - 支持多种OSS
func GetPrivateURLExact(domain, key, accessKey, secretKey string) (string, error) {
	fmt.Printf("domain: %s, key: %s, accessKey: %s, secretKey: %s\n", domain, key, accessKey, secretKey)

	// 判断是否为S3兼容域名（如七牛云S3、AWS S3、阿里云OSS等）
	if isS3CompatibleDomain(domain) {
		return generateS3PresignedURL(domain, key, accessKey, secretKey)
	}

	// 七牛云传统域名处理
	return generateQiniuTraditionalURL(domain, key, accessKey, secretKey)
}

// generateQiniuTraditionalURL 生成七牛云传统私有URL
func generateQiniuTraditionalURL(domain, key, accessKey, secretKey string) (string, error) {
	mac := auth.New(accessKey, secretKey)
	deadline := time.Now().Add(time.Second * 3600).Unix()
	privateURL := storage.MakePrivateURLv2(mac, domain, key, deadline)
	fmt.Printf("七牛云传统URL: %s\n", privateURL)
	return privateURL, nil
}

// isS3CompatibleDomain 判断是否为S3兼容域名
func isS3CompatibleDomain(domain string) bool {
	// 如果域名包含".s3."，很可能是S3兼容接口
	return strings.Contains(domain, ".s3.")
}

// generateS3PresignedURL 生成S3兼容的预签名URL
func generateS3PresignedURL(domain, key, accessKey, secretKey string) (string, error) {
	// 从域名中提取bucket和endpoint
	bucket, endpoint, region := extractS3Info(domain)

	fmt.Printf("S3模式 - bucket: %s, endpoint: %s, region: %s\n", bucket, endpoint, region)

	// 创建AWS S3配置
	config := &aws.Config{
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false), // 使用虚拟主机样式
		DisableSSL:       aws.Bool(false),
	}

	// 创建session
	sess, err := session.NewSession(config)
	if err != nil {
		return "", fmt.Errorf("创建S3 session失败: %v", err)
	}

	// 创建S3客户端
	svc := s3.New(sess)

	// 生成预签名URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(3600 * time.Second) // 1小时有效
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}

	fmt.Printf("S3预签名URL: %s\n", url)
	return url, nil
}

// extractS3Info 从S3兼容域名中提取信息
func extractS3Info(domain string) (bucket, endpoint, region string) {
	// 七牛云S3格式：bucket.s3.region.qiniucs.com
	if strings.Contains(domain, ".qiniucs.com") {
		parts := strings.Split(domain, ".")
		if len(parts) >= 4 {
			bucket = parts[0]
			region = parts[2]
			endpoint = fmt.Sprintf("https://s3.%s.qiniucs.com", region)
			return
		}
	}

	// AWS S3格式：bucket.s3.region.amazonaws.com
	if strings.Contains(domain, ".amazonaws.com") {
		parts := strings.Split(domain, ".")
		if len(parts) >= 5 {
			bucket = parts[0]
			region = parts[2]
			endpoint = fmt.Sprintf("https://s3.%s.amazonaws.com", region)
			return
		}
	}

	// 默认处理：尝试从常见格式中提取
	parts := strings.Split(domain, ".")
	if len(parts) > 0 {
		bucket = parts[0]
		region = "us-east-1" // 默认区域
		endpoint = fmt.Sprintf("https://%s", strings.Join(parts[1:], "."))
	}

	return
}
