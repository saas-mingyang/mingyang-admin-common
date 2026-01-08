package cloudfile

import (
	"context"

	"mingyang.com/admin-simple-admin-file/ent/cloudfile"

	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"
	"mingyang.com/admin-simple-admin-file/internal/utils/dberrorhandler"

	"github.com/saas-mingyang/mingyang-admin-common/i18n"

	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCloudFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileByIdLogic {
	return &GetCloudFileByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCloudFileByIdLogic) GetCloudFileById(req *types.BaseIDInfo) (*types.CloudFileInfoResp, error) {
	data, err := l.svcCtx.DB.CloudFile.Query().Where(cloudfile.IDEQ(*req.Id)).WithStorageProviders().
		First(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.CloudFileInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  l.svcCtx.Trans.Trans(l.ctx, i18n.Success),
		},
		Data: types.CloudFileInfo{
			BaseIDInfo: types.BaseIDInfo{
				Id:        &data.ID,
				CreatedAt: pointy.GetPointer(data.CreatedAt.UnixMilli()),
				UpdatedAt: pointy.GetPointer(data.UpdatedAt.UnixMilli()),
			},
			State:      &data.State,
			Name:       &data.Name,
			Url:        &data.URL,
			Size:       &data.Size,
			FileType:   &data.FileType,
			UserId:     &data.UserID,
			ProviderId: &data.Edges.StorageProviders.ID,
		},
	}, nil
}
