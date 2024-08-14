package storageprovider

import (
	"context"
	"github.com/suyuan32/simple-admin-common/orm/ent/entctx/tenantctx"
	"github.com/suyuan32/simple-admin-file/ent/cloudfile"
	"github.com/suyuan32/simple-admin-file/ent/storageprovider"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/cloud"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteStorageProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteStorageProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteStorageProviderLogic {
	return &DeleteStorageProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteStorageProviderLogic) DeleteStorageProvider(req *types.IDsReq) (*types.BaseMsgResp, error) {
	tenantId := tenantctx.GetTenantIDFromCtx(l.ctx)

	check, err := l.svcCtx.DB.CloudFile.Query().Where(cloudfile.HasStorageProvidersWith(storageprovider.IDIn(req.Ids...))).
		Count(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	if check != 0 {
		return nil, errorx.NewCodeInvalidArgumentError("storage_provider.hasFileError")
	}

	_, err = l.svcCtx.DB.StorageProvider.Delete().Where(storageprovider.IDIn(req.Ids...)).Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	err = cloud.AddTenantCloudServiceGroup(l.svcCtx.DB, l.svcCtx.CloudStorage, tenantId)
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess)}, nil
}
