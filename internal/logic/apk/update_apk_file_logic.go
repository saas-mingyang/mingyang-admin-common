package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"
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

func (l *UpdateApkFileLogic) UpdateApkFile(req *types.ApkUpdateReq) (resp *types.BaseMsgResp, err error) {
	err = l.svcCtx.DB.Apk.UpdateOneID(*req.Id).
		SetDescription(req.Description).
		SetPackageName(req.PackageName).
		SetUpdateLog(req.UpdateLog).
		Exec(l.ctx)
	if err != nil {
		return nil, errorx.NewCodeInternalError(err.Error())
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.UpdateSuccess)}, nil
}
