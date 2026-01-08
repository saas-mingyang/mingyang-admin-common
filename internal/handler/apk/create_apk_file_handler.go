package apk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang.com/admin-simple-admin-file/internal/logic/apk"
	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"
)

// swagger:route post /apk/create apk CreateApkFile
//
// Create apk file information | 创建apk文件
//
// Create apk file information | 创建apk文件
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ApkInfo
//
// Responses:
//  200: BaseMsgResp

func CreateApkFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApkInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apk.NewCreateApkFileLogic(r.Context(), svcCtx)
		resp, err := l.CreateApkFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
