package storageprovider

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/tenantctx"
	"github.com/suyuan32/simple-admin-file-tenant/internal/svc"
	"github.com/suyuan32/simple-admin-file-tenant/internal/types"
	"github.com/suyuan32/simple-admin-file-tenant/internal/utils/cloud"
	"github.com/suyuan32/simple-admin-file-tenant/internal/utils/dberrorhandler"

	"github.com/saas-mingyang/mingyang-admin-common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStorageProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStorageProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStorageProviderLogic {
	return &UpdateStorageProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStorageProviderLogic) UpdateStorageProvider(req *types.StorageProviderInfo) (*types.BaseMsgResp, error) {
	tenantId := tenantctx.GetTenantIDFromCtx(l.ctx)

	err := l.svcCtx.DB.StorageProvider.UpdateOneID(*req.Id).
		SetNotNilState(req.State).
		SetNotNilName(req.Name).
		SetNotNilBucket(req.Bucket).
		SetNotNilSecretID(req.SecretId).
		SetNotNilSecretKey(req.SecretKey).
		SetNotNilRegion(req.Region).
		SetNotNilIsDefault(req.IsDefault).
		SetNotNilFolder(req.Folder).
		SetNotNilEndpoint(req.Endpoint).
		SetNotNilUseCdn(req.UseCdn).
		SetNotNilCdnURL(req.CdnUrl).
		Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	err = cloud.AddTenantCloudServiceGroup(l.svcCtx.DB, l.svcCtx.CloudStorage, tenantId)
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.UpdateSuccess)}, nil
}
