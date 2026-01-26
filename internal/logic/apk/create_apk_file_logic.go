package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/utils/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
	ent2 "mingyang.com/admin-simple-admin-file/ent"
	"mingyang.com/admin-simple-admin-file/ent/apk"
	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"
	"mingyang.com/admin-simple-admin-file/internal/utils/dberrorhandler"
	"mingyang.com/admin-simple-admin-file/internal/utils/entx"
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

	err = entx.WithTx(l.ctx, l.svcCtx.DB, func(tx *ent2.Tx) error {
		// 1. 检查同名同版本是否已存在
		existingByName, err := tx.Apk.Query().
			Where(apk.NameEQ(req.Name), apk.VersionEQ(req.Version)).
			First(l.ctx)

		if err != nil && !ent2.IsNotFound(err) {
			l.Logger.Errorf("查询应用错误: %v", err)
			return err
		}

		// 2. 如果存在同名同版本，则更新（不是删除）
		if err == nil && existingByName != nil {
			_, err = tx.Apk.UpdateOneID(existingByName.ID).
				SetVersionCode(req.VersionCode).
				SetFileURL(req.AppStoreUrl).
				SetNillableDescription(req.Description).
				SetNillableUpdateLog(req.UpdateLog).
				SetNillableIsForceUpdate(req.IsForceUpdate).
				SetPackageName(req.PackageName).
				SetNillableFileSize(req.FileSize).
				SetCategory(req.Category).
				SetFileID(req.FileInfo.FileId).
				Save(l.ctx)
			return err
		} else {
			// 3. 创建新记录
			_, err = tx.Apk.Create().
				SetID(sonyflake.NextID()).
				SetName(req.Name).
				SetVersion(req.Version).
				SetVersionCode(req.VersionCode). // 使用用户填写的版本号
				SetFileURL(req.AppStoreUrl).
				SetNillableDescription(req.Description).
				SetNillableUpdateLog(req.UpdateLog).
				SetNillableIsForceUpdate(req.IsForceUpdate).
				SetPackageName(req.PackageName).
				SetNillableFileSize(req.FileSize).
				SetCategory(req.Category).
				SetFileID(req.FileInfo.FileId).
				Save(l.ctx)
		}

		if err != nil {
			l.Logger.Errorf("创建应用错误: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.CreateSuccess)}, nil
}
