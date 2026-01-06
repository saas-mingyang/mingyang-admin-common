package cloudfile

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mingyang-admin-simple-admin-file/internal/logic/cloudfile"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
)

// swagger:route post /cloud_file/upload_progress cloudfile GetCloudFileUploadProgress
//
// Get cloud file  upload progress | 获取云文件上传进度
//
// Get cloud file  upload progress | 获取云文件上传进度
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: UUIDReq
//
// Responses:
//  200: UploadProgressResp

func GetCloudFileUploadProgressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UUIDReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cloudfile.NewGetCloudFileUploadProgressLogic(r.Context(), svcCtx)
		resp, err := l.GetCloudFileUploadProgress(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
