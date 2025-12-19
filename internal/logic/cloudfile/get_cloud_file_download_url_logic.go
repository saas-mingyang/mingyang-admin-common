package cloudfile

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang-admin-simple-admin-file/ent/cloudfile"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
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

func (l *GetCloudFileDownloadUrlLogic) GetCloudFileDownloadUrl(req *types.BaseIDInfo) (resp *types.CloudFileInfoResp, err error) {
	file, err := l.svcCtx.DB.CloudFile.
		Query().
		Where(cloudfile.ID(*req.Id)).
		WithStorageProviders().
		First(l.ctx)
	if err != nil || file == nil {
		return nil, errorx.NewCodeInvalidArgumentError(err.Error())
	}
	providers := file.Edges.StorageProviders
	if providers == nil {
		return nil, errorx.NewCodeInvalidArgumentError(err.Error())
	}
	privateURL, err := GetPrivateURLExact(providers.CdnURL, file.URL, providers.SecretID, providers.SecretKey)
	if err != nil {
		return nil, errorx.NewCodeInternalError(err.Error())
	}

	resp = &types.CloudFileInfoResp{}

	id := file.ID
	resp.Data = types.CloudFileInfo{
		Url:        &privateURL,
		Size:       &file.Size,
		Name:       &file.Name,
		FileType:   &file.FileType,
		UserId:     &file.UserID,
		ProviderId: &providers.ID,
		BaseIDInfo: types.BaseIDInfo{
			Id: &id,
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
