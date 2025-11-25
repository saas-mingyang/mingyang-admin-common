package apk

import (
	"context"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApkFileLogic {
	return &UpdateApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApkFileLogic) UpdateApkFile(req *types.ApkInfo) (resp *types.BaseMsgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
