package cloud

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/tenantctx"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entenum"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang-admin-simple-admin-file/ent"
	"mingyang-admin-simple-admin-file/ent/storageprovider"
)

type CloudServiceGroup struct {
	Service map[uint64]*CloudService
}

type CloudService struct {
	CloudStorage map[string]*s3.S3

	ProviderData map[string]struct {
		Id       uint64
		Folder   string
		Bucket   string
		Endpoint string
		UseCdn   bool
		CdnUrl   string
	}

	DefaultProvider string
}

// NewCloudServiceGroup returns the S3 service client group
func NewCloudServiceGroup(db *ent.Client) *CloudServiceGroup {
	cloudServices := &CloudServiceGroup{}
	cloudServices.Service = make(map[uint64]*CloudService)

	_ = AddTenantCloudServiceGroup(db, cloudServices, entenum.TenantDefaultId)

	return cloudServices
}

func AddTenantCloudServiceGroup(db *ent.Client, service *CloudServiceGroup, tenantId uint64) error {
	data, err := db.StorageProvider.Query().Where(storageprovider.StateEQ(true), storageprovider.TenantIDEQ(tenantId)).All(tenantctx.AdminCtx(context.Background()))
	if err != nil {
		logx.Errorw("failed to load provider config from database, make sure database has been initialize and has config data",
			logx.Field("detail", err))
		return err
	}

	service.Service[tenantId] = &CloudService{CloudStorage: make(map[string]*s3.S3), ProviderData: make(map[string]struct {
		Id       uint64
		Folder   string
		Bucket   string
		Endpoint string
		UseCdn   bool
		CdnUrl   string
	}), DefaultProvider: ""}

	for _, v := range data {
		sess := session.Must(session.NewSession(
			&aws.Config{
				Region:      aws.String(v.Region),
				Credentials: credentials.NewStaticCredentials(v.SecretID, v.SecretKey, ""),
				Endpoint:    aws.String(v.Endpoint),
			},
		))
		svc := s3.New(sess)

		service.Service[tenantId].CloudStorage[v.Name] = svc
		service.Service[tenantId].ProviderData[v.Name] = struct {
			Id       uint64
			Folder   string
			Bucket   string
			Endpoint string
			UseCdn   bool
			CdnUrl   string
		}{Id: v.ID, Folder: v.Folder, Bucket: v.Bucket, Endpoint: v.Endpoint, UseCdn: v.UseCdn, CdnUrl: v.CdnURL}

		if v.IsDefault {
			service.Service[tenantId].DefaultProvider = v.Name
		}
	}

	return err
}
