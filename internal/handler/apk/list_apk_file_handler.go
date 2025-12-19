package apk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang-admin-simple-admin-file/internal/logic/apk"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
)

// swagger:route post /apk/list apk ListApkFile
//
// List apk file information | 获取apk文件列表
//
// List apk file information | 获取apk文件列表
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ApkFileListReq
//
// Responses:
//  200: ApkFileListResp

func ListApkFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApkFileListReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apk.NewListApkFileLogic(r.Context(), svcCtx)
		resp, err := l.ListApkFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
