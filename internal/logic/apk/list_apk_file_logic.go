package apk

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/saas-mingyang/mingyang-admin-common/utils/pointy"
	"mingyang.com/admin-simple-admin-file/ent/apk"
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
	data, err := l.svcCtx.DB.Apk.Query().Where(predicates...).
		Page(l.ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}
	resp = &types.ApkFileListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.PageDetails.Total
	for _, v := range data.List {
		resp.Data.Data = append(resp.Data.Data,
			types.ApkInfo{
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
			})
	}
	return resp, nil
}
