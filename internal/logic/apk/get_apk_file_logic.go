package apk

import (
	"context"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApkFileLogic {
	return &GetApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApkFileLogic) GetApkFile(req *types.IdsReq) (resp *types.BaseMsgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
