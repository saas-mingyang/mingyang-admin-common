package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"mingyang-admin-simple-admin-file/ent/apk"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"

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
	_, err = l.svcCtx.DB.Apk.Delete().Where(apk.IDIn(req.Ids...)).Exec(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess)}, nil
}
