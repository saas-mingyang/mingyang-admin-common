package cloudfile

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/utils/sonyflake"

	"github.com/zeromicro/go-zero/core/errorx"

	"mingyang-admin-simple-admin-file/ent"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"

	"github.com/saas-mingyang/mingyang-admin-common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCloudFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCloudFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCloudFileLogic {
	return &CreateCloudFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCloudFileLogic) CreateCloudFile(req *types.CloudFileInfo) (*types.BaseMsgResp, error) {
	// check storage provider exist
	_, err := l.svcCtx.DB.StorageProvider.Get(l.ctx, *req.ProviderId)
	switch {
	case ent.IsNotFound(err):
		return nil, errorx.NewCodeInvalidArgumentError("storage_provider.StorageProviderNotExist")
	case err != nil:
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	query := l.svcCtx.DB.CloudFile.Create().
		SetNotNilState(req.State).
		SetNotNilName(req.Name).
		SetNotNilURL(req.Url).
		SetID(sonyflake.NextID()).
		SetNotNilSize(req.Size).
		SetNotNilFileType(req.FileType).
		SetNotNilUserID(req.UserId)

	if req.ProviderId != nil {
		query = query.SetStorageProvidersID(*req.ProviderId)
	}

	if req.TagIds != nil {
		query = query.AddTagIDs(req.TagIds...)
	}

	_, err = query.Save(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.CreateSuccess)}, nil
}
