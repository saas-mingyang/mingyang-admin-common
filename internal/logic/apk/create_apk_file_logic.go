package apk

import (
	"context"
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/utils/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang-admin-simple-admin-file/ent/apk"
	"mingyang-admin-simple-admin-file/internal/logic/cloudfile"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"
	"strings"
)

type CreateApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApkFileLogic {
	return &CreateApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateApkFileLogic) CreateApkFile(req *types.ApkInfo) (resp *types.BaseMsgResp, err error) {
	var builder strings.Builder
	builder.WriteString(req.Name)
	builder.WriteString("_")
	builder.WriteString(req.Version)
	existing, err := l.svcCtx.DB.Apk.Query().
		Where(apk.VersionCode(builder.String())).
		First(l.ctx)
	if existing != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	downloadUrlLogic := cloudfile.NewGetCloudFileDownloadUrlLogic(l.ctx, l.svcCtx)
	result, err := downloadUrlLogic.GetCloudFileDownloadUrl(&types.BaseIDInfo{Id: req.FileId})
	if err != nil {
		fmt.Printf("get cloud file download url error: %v", err)
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	_, err = l.svcCtx.DB.Apk.Create().
		SetID(sonyflake.NextID()).
		SetName(req.Name).
		SetVersion(req.Version).
		SetVersionCode(builder.String()).
		SetFileID(*req.FileId).
		SetFilePath(req.FilePath).
		SetDescription(req.Description).
		SetUpdateLog(req.UpdateLog).
		SetIsForceUpdate(req.IsForceUpdate).
		SetPackageName(req.PackageName).
		SetFileSize(*result.Data.Size).
		Save(l.ctx)
	if err != nil {
		fmt.Printf("create apk error: %v", err)
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.CreateSuccess)}, nil
}
