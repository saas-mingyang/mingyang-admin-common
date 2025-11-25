package apk

import (
	"context"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApkFileLogic {
	return &ListApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApkFileLogic) ListApkFile(req *types.ApkFileListReq) (resp *types.ApkFileListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
