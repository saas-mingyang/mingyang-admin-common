package cloudfile

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"github.com/saas-mingyang/mingyang-admin-common/utils/uuidx"
	"github.com/suyuan32/simple-admin-file-tenant/ent/cloudfile"
	"github.com/suyuan32/simple-admin-file-tenant/internal/svc"
	"github.com/suyuan32/simple-admin-file-tenant/internal/types"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type GetCloudFileDownloadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCloudFileDownloadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileDownloadUrlLogic {
	return &GetCloudFileDownloadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCloudFileDownloadUrlLogic) GetCloudFileDownloadUrl(req *types.UUIDReq) (resp *types.CloudFileInfoResp, err error) {
	file, err := l.svcCtx.DB.CloudFile.
		Query().
		Where(cloudfile.ID(uuidx.ParseUUIDString(req.Id))).
		WithStorageProviders().
		First(l.ctx)
	if err != nil || file == nil {
		return nil, errorx.NewCodeInvalidArgumentError("cloud_file.CloudFileNotExist")
	}
	providers := file.Edges.StorageProviders
	if providers == nil {
		return nil, errorx.NewCodeInvalidArgumentError("cloud_file.CloudFileNotExist")
	}
	privateURL, err := GetPrivateURLExact(providers.CdnURL, file.URL, providers.SecretID, providers.SecretKey)
	if err != nil {
		return nil, errorx.NewCodeInternalError(err.Error())
	}

	resp = &types.CloudFileInfoResp{}

	id := file.ID
	resp.Data = types.CloudFileInfo{
		Url: &privateURL,
		BaseUUIDInfo: types.BaseUUIDInfo{
			Id: pointy.GetPointer(id.String()),
		},
	}

	return resp, nil
}

func GetPrivateURLExact(domain, key, accessKey, secretKey string) (string, error) {
	fmt.Printf("domain: %s, key: %s, accessKey: %s, secretKey: %s\n", domain, key, accessKey, secretKey)
	mac := auth.New(accessKey, secretKey)
	deadline := time.Now().Add(time.Second * 3600).Unix()
	privateURL := storage.MakePrivateURLv2(mac, domain, key, deadline)
	return privateURL, nil
}
