package cloudfile

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"mingyang.com/admin-simple-admin-file/internal/logic/cloudfile"
	"mingyang.com/admin-simple-admin-file/internal/svc"
	"net/http"
	"time"
)

// swagger:route post /cloud_file/upload cloudfile Upload
//
// Cloud file upload | 上传文件
//
// Cloud file upload | 上传文件
//
// Responses:
//  200: CloudFileInfoResp

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 创建一个新的上下文，设置30分钟超时
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
		defer cancel() // 确保函数退出时取消上下文

		// 使用新的上下文创建请求
		r = r.WithContext(ctx)

		l := cloudfile.NewUploadLogic(r, svcCtx)
		resp, err := l.Upload()
		if err != nil {
			err = svcCtx.Trans.TransError(ctx, err)
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
