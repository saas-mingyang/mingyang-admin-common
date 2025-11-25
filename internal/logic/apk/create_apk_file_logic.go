package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/utils/sonyflake"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang-admin-simple-admin-core/rpc/ent"
	"mingyang-admin-simple-admin-file/ent/apk"
	"mingyang-admin-simple-admin-file/internal/logic/cloudfile"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
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
	if err != nil && !ent.IsNotFound(err) {
		return nil, errorx.NewCodeInternalError("查询版本号失败: " + err.Error())
	}
	if existing != nil {
		return nil, errorx.NewCodeError(400, "版本号 "+req.VersionCode+" 已存在")
	}
	downloadUrlLogic := cloudfile.NewGetCloudFileDownloadUrlLogic(l.ctx, l.svcCtx)
	result, err := downloadUrlLogic.GetCloudFileDownloadUrl(&types.UUIDReq{Id: req.FileId})
	if err != nil {
		return nil, errorx.NewCodeInternalError(err.Error())
	}
	data := result.Data
	_, err = l.svcCtx.DB.Apk.Create().
		SetID(sonyflake.NextID()).
		SetName(req.Name).
		SetVersion(req.Version).
		SetVersionCode(builder.String()).
		SetFileID(req.FileId).
		SetFilePath(*data.Url).
		SetDescription(req.Description).
		SetUpdateLog(req.UpdateLog).
		SetIsForceUpdate(req.IsForceUpdate).
		SetMd5(req.Md5).
		SetSha1(req.Sha1).
		SetSha256(req.Sha256).
		SetPackageName(req.PackageName).
		SetFileSize(*data.Size).
		Save(l.ctx)
	if err != nil {
		return nil, errorx.NewCodeInternalError(err.Error())
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.CreateSuccess)}, nil
}
