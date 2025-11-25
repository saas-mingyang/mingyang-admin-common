package apk

import (
	"context"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadApkFileLogic {
	return &DownloadApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadApkFileLogic) DownloadApkFile(req *types.IdsReq) (resp *types.BaseMsgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
