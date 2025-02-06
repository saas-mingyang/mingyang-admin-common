package cloudfile

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/suyuan32/simple-admin-common/orm/ent/entctx/tenantctx"
	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/suyuan32/simple-admin-file/internal/utils/cloud"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-file/ent/cloudfile"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCloudFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCloudFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCloudFileLogic {
	return &DeleteCloudFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCloudFileLogic) DeleteCloudFile(req *types.UUIDsReq) (*types.BaseMsgResp, error) {
	if l.svcCtx.Config.UploadConf.DeleteFileWithCloud {
		tenantId := tenantctx.GetTenantIDFromCtx(l.ctx)
		if _, ok := l.svcCtx.CloudStorage.Service[tenantId]; !ok {
			err := cloud.AddTenantCloudServiceGroup(l.svcCtx.DB, l.svcCtx.CloudStorage, tenantId)
			if err != nil {
				if ent.IsNotFound(err) {
					return nil, errorx.NewCodeInternalError("storage_provider.StorageProviderNotExist")
				}
				return nil, errorx.NewCodeInternalError("storage_provider.failedLoadProviderConfig")
			}
		}

		data, err := l.svcCtx.DB.CloudFile.Query().Where(cloudfile.IDIn(uuidx.ParseUUIDSlice(req.Ids)...)).
			WithStorageProviders().All(l.ctx)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
		}

		for _, v := range data {
			if client, ok := l.svcCtx.CloudStorage.Service[tenantId].CloudStorage[v.Edges.StorageProviders.Name]; ok {
				if !l.svcCtx.CloudStorage.Service[tenantId].ProviderData[v.Edges.StorageProviders.Name].UseCdn {
					keyData := strings.Split(v.URL, *l.svcCtx.CloudStorage.Service[tenantId].CloudStorage[v.Edges.StorageProviders.Name].Config.Endpoint)
					if len(keyData) != 2 {
						logx.Errorw("failed to find the key of the cloud file", logx.Field("data", req))
						return nil, errorx.NewCodeInternalError(i18n.Failed)
					}
					_, err = client.DeleteObject(&s3.DeleteObjectInput{
						Bucket: aws.String(l.svcCtx.CloudStorage.Service[tenantId].ProviderData[v.Edges.StorageProviders.Name].Bucket),
						Key:    aws.String(keyData[1]),
					})
				} else {
					keyData := strings.TrimPrefix(v.URL, l.svcCtx.CloudStorage.Service[tenantId].ProviderData[v.Edges.StorageProviders.Name].CdnUrl)
					if len(keyData) < 2 {
						logx.Errorw("failed to find the key of the cloud file", logx.Field("data", req))
						return nil, errorx.NewCodeInternalError(i18n.Failed)
					}
					_, err = client.DeleteObject(&s3.DeleteObjectInput{
						Bucket: aws.String(l.svcCtx.CloudStorage.Service[tenantId].ProviderData[v.Edges.StorageProviders.Name].Bucket),
						Key:    aws.String(keyData),
					})
				}
				if err != nil {
					logx.Errorw("failed to delete the cloud file", logx.Field("detail", err), logx.Field("data", req))
					return nil, errorx.NewCodeInternalError(i18n.Failed)
				}
			}
		}
	}

	_, err := l.svcCtx.DB.CloudFile.Delete().Where(cloudfile.IDIn(uuidx.ParseUUIDSlice(req.Ids)...)).Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess)}, nil
}
