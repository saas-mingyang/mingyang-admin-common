package apk

import (
	"context"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApkFileLogic {
	return &DeleteApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApkFileLogic) DeleteApkFile(req *types.IdsReq) (resp *types.BaseMsgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
