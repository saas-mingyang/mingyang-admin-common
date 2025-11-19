package cloudfile

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"github.com/saas-mingyang/mingyang-admin-common/utils/uuidx"
	"github.com/suyuan32/simple-admin-file-tenant/ent/cloudfile"
	"github.com/zeromicro/go-zero/core/errorx"
	"time"

	"github.com/suyuan32/simple-admin-file-tenant/internal/svc"
	"github.com/suyuan32/simple-admin-file-tenant/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	privateURL := GetPrivateURL(providers.CdnURL, file.URL, providers.SecretID, providers.SecretKey)

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

// GetPrivateURL 获取私有下载链接
func GetPrivateURL(domain, key, accessKey, secretKey string) string {
	fmt.Printf("domain: %s, key: %s, accessKey: %s, secretKey: %s", domain, key, accessKey, secretKey)
	mac := auth.New(accessKey, secretKey)
	deadline := time.Now().Unix() + 3600
	privateURL := storage.MakePrivateURL(mac, domain, key, deadline)
	return privateURL
}
