package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"mingyang.com/admin-simple-admin-file/ent/apk"
	"mingyang.com/admin-simple-admin-file/ent/cloudfile"
	"mingyang.com/admin-simple-admin-file/ent/predicate"
	"mingyang.com/admin-simple-admin-file/internal/utils/dberrorhandler"

	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApkFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApkFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApkFileLogic {
	return &ListApkFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApkFileLogic) ListApkFile(req *types.ApkFileListReq) (resp *types.ApkFileListResp, err error) {
	var predicates []predicate.Apk
	if req.Name != nil {
		predicates = append(predicates, apk.NameContains(*req.Name))
	}
	if req.VersionCode != nil {
		predicates = append(predicates, apk.VersionCodeContains(*req.VersionCode))
	}
	if req.Version != nil {
		predicates = append(predicates, apk.VersionContains(*req.Version))
	}
	if req.PackageName != nil {
		predicates = append(predicates, apk.PackageNameContains(*req.PackageName))
	}
	if req.Category != nil {
		predicates = append(predicates, apk.CategoryEQ(*req.Category))
	}
	if req.Status != nil {
		predicates = append(predicates, apk.StatusEQ(*req.Status))
	}
	data, err := l.svcCtx.DB.Apk.Query().Where(predicates...).
		Page(l.ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	resp = &types.ApkFileListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.PageDetails.Total
	for _, v := range data.List {
		apkInfo := types.ApkInfo{
			BaseIDInfo: types.BaseIDInfo{
				Id:        &v.ID,
				CreatedAt: pointy.GetPointer(v.CreatedAt.UnixMilli()),
				UpdatedAt: pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			},

			Name:        v.Name,
			Version:     v.Version,
			VersionCode: v.VersionCode,
			AppStoreUrl: v.FileURL,
			PackageName: v.PackageName,
			Description: &v.Description,
			Category:    v.Category,
			FileId:      v.FileID,
			Status:      v.Status,
			FileSize:    pointy.GetPointer(v.FileSize),
		}
		if apkInfo.FileId != 0 {
			file, err := l.svcCtx.DB.CloudFile.Query().Where(cloudfile.ID(apkInfo.FileId)).First(l.ctx)
			if err != nil {
				l.Error("query file error: %v", err)
			}
			if file != nil {
				apkInfo.FileSize = pointy.GetPointer(file.Size)
				apkInfo.FileName = file.Name
			}
		}
		resp.Data.Data = append(resp.Data.Data,
			apkInfo,
		)
	}
	return resp, nil
}
