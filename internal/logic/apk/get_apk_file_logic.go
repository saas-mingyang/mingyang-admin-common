package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"mingyang-admin-simple-admin-file/ent/apk"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApkFileLogic {
	return &GetApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApkFileLogic) GetApkFile(req *types.IDReq) (resp *types.ApkFileInfoResp, err error) {
	data, err := l.svcCtx.DB.Apk.Query().Where(apk.IDEQ(req.Id)).First(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.ApkFileInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  l.svcCtx.Trans.Trans(l.ctx, i18n.Success),
		},
		Data: types.ApkInfo{
			BaseIDInfo: types.BaseIDInfo{
				Id:        &data.ID,
				CreatedAt: pointy.GetPointer(data.CreatedAt.UnixMilli()),
				UpdatedAt: pointy.GetPointer(data.UpdatedAt.UnixMilli()),
			},
			FileId:        &data.FileID,
			Name:          data.Name,
			FilePath:      data.FilePath,
			FileSize:      data.FileSize,
			Version:       data.Version,
			VersionCode:   data.VersionCode,
			PackageName:   data.PackageName,
			Description:   data.Description,
			UpdateLog:     data.UpdateLog,
			IsForceUpdate: data.IsForceUpdate,
			DownloadCount: data.DownloadCount,
		},
	}, nil
}
