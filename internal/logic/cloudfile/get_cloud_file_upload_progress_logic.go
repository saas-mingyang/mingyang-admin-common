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

func (l *GetCloudFileUploadProgressLogic) GetCloudFileUploadProgress(req *types.UUIDReq) (resp *types.UploadProgressResp, err error) {

	return
}
