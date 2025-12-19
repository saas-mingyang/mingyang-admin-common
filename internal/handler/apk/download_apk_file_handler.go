package apk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang-admin-simple-admin-file/internal/logic/apk"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
)

// swagger:route post /apk/download_url apk DownloadApkFile
//
// Download apk file information | 下载apk文件
//
// Download apk file information | 下载apk文件
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: IDReq
//
// Responses:
//  200: CloudFileInfoResp

func DownloadApkFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IDReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apk.NewDownloadApkFileLogic(r.Context(), svcCtx)
		resp, err := l.DownloadApkFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
