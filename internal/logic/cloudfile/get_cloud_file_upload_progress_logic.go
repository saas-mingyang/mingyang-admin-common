package cloudfile

import (
	"context"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

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
func (l *GetCloudFileUploadProgressLogic) GetCloudFileUploadProgress(req *types.UUIDReq) (resp *types.UploadProgressResp, err error) {
	logx.Infow("查询多个上传进度",
		logx.Field("UploadIDs", req.Id),
		logx.Field("userId", l.ctx.Value("userId")))

	return
}
