package cloudfile

import (
	"context"

	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileUploadProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCloudFileUploadProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileUploadProgressLogic {
	return &GetCloudFileUploadProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCloudFileUploadProgress  查询单个上传进度
// 参数:
//   - req: 上传进度查询请求，包含uploadId字段
//
// 返回:
//   - resp: 上传进度查询响应，包含进度详情或错误信息
//   - err: 错误信息（如果有）
//
// 说明:
//   - 根据uploadId查询特定上传任务的进度信息
//   - 支持断点续传、进度显示等场景
func (l *GetCloudFileUploadProgressLogic) GetCloudFileUploadProgress(req *types.IDReq) (resp *types.UploadProgressResp, err error) {
	progress, err := GetUploadProgress(req.Id)

	if err != nil {
		return nil, err
	}
	return &types.UploadProgressResp{
		UploadID:    req.Id,
		FileName:    progress.FileName,
		Key:         progress.Key,
		Bucket:      progress.Bucket,
		Provider:    progress.Provider,
		Status:      progress.Status,
		TotalSize:   progress.TotalSize,
		TotalParts:  progress.TotalParts,
		CurrentPart: progress.CurrentPart,
		Uploaded:    progress.Uploaded,
		Percentage:  progress.Percentage,
		Speed:       progress.Speed,
		StartTime:   progress.StartTime,
	}, nil
}
