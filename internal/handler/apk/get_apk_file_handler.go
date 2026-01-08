package apk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang.com/admin-simple-admin-file/internal/logic/apk"
	"mingyang.com/admin-simple-admin-file/internal/svc"
	"mingyang.com/admin-simple-admin-file/internal/types"
)

// swagger:route post /apk/get apk GetApkFile
//
// Get apk file information | 获取apk文件
//
// Get apk file information | 获取apk文件
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: IdsReq
//
// Responses:
//  200: BaseMsgResp

func GetApkFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IDReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apk.NewGetApkFileLogic(r.Context(), svcCtx)
		resp, err := l.GetApkFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
