package apk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang-admin-simple-admin-file/internal/logic/apk"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
)

// swagger:route post /apk/delete apk DeleteApkFile
//
// Delete apk file information | 删除apk文件
//
// Delete apk file information | 删除apk文件
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: IdsReq
//
// Responses:
//  200: BaseMsgResp

func DeleteApkFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdsReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apk.NewDeleteApkFileLogic(r.Context(), svcCtx)
		resp, err := l.DeleteApkFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
