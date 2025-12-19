package apk

import (
	"context"
	"github.com/zeromicro/go-zero/core/errorx"
	"mingyang-admin-simple-admin-file/internal/logic/cloudfile"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"
	"strconv"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadApkFileLogic {
	return &DownloadApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadApkFileLogic) DownloadApkFile(req *types.IDReq) (resp *types.CloudFileInfoResp, err error) {
	file, err := l.svcCtx.DB.Apk.Get(l.ctx, req.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	fileId, err := strconv.ParseUint(file.FileURL, 10, 64)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(l.svcCtx.Trans.Trans(l.ctx, "file.fileCategoryIsAndroid"))
	}
	downloadUrlLogic := cloudfile.NewGetCloudFileDownloadUrlLogic(l.ctx, l.svcCtx)

	result, err := downloadUrlLogic.GetCloudFileDownloadUrl(&types.BaseIDInfo{Id: &fileId})
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	//对下载次数+1
	_, err = l.svcCtx.DB.Apk.UpdateOneID(req.Id).SetDownloadCount(file.DownloadCount + 1).Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	return result, nil
}
