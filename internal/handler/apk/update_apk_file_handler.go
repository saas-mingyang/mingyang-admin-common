package apk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang-admin-simple-admin-file/internal/logic/apk"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
)

// swagger:route post /apk/update apk UpdateApkFile
//
// update apk file information | 创建apk文件
//
// update apk file information | 创建apk文件
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ApkInfo
//
// Responses:
//  200: BaseMsgResp

func UpdateApkFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApkUpdateReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apk.NewUpdateApkFileLogic(r.Context(), svcCtx)
		resp, err := l.UpdateApkFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
