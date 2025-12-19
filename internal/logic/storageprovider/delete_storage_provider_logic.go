package storageprovider

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/tenantctx"
	"github.com/saas-mingyang/mingyang-admin-common/utils/convert"
	"github.com/zeromicro/go-zero/core/errorx"
	"mingyang-admin-simple-admin-file/ent/cloudfile"
	"mingyang-admin-simple-admin-file/ent/storageprovider"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
	"mingyang-admin-simple-admin-file/internal/utils/cloud"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"

	"github.com/saas-mingyang/mingyang-admin-common/i18n"
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

	check, err := l.svcCtx.DB.CloudFile.Query().Where(cloudfile.HasStorageProvidersWith(storageprovider.IDIn(convert.StringSliceToUint64Slice(req.Ids)...))).
		Count(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	if check != 0 {
		return nil, errorx.NewCodeInvalidArgumentError("storage_provider.hasFileError")
	}

	_, err = l.svcCtx.DB.StorageProvider.Delete().Where(storageprovider.IDIn(convert.StringSliceToUint64Slice(req.Ids)...)).Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	err = cloud.AddTenantCloudServiceGroup(l.svcCtx.DB, l.svcCtx.CloudStorage, tenantId)
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess)}, nil
}
