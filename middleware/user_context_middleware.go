package middleware

import (
	"context"
	"net/http"

	"github.com/saas-mingyang/mingyang-admin-common/enum/common"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/tenantctx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/enum"
	"google.golang.org/grpc/metadata"
)

// UserContextMiddleware 将网关通过 HTTP Header 透传过来的用户身份信息
// (X-User-ID / X-Role-ID / X-Tenant-ID) 重新写回当前请求的 context,
// 使 userctx.GetUserIDFromCtx / tenantctx.GetTenantIDFromCtx 能正常取值,
// 同时追加到 outgoing gRPC metadata, 供下游 RPC 使用。
//
// 适用于位于网关之后的各业务 API 服务, 通过 server.Use(...) 全局注册。
type UserContextMiddleware struct{}

func NewUserContextMiddleware() *UserContextMiddleware {
	return &UserContextMiddleware{}
}

func (m *UserContextMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if userId := r.Header.Get(common.HeaderXUserID); userId != common.EmptyString {
			ctx = context.WithValue(ctx, common.CtxKeyUserID, userId)
			ctx = metadata.AppendToOutgoingContext(ctx, enum.UserIdRpcCtxKey, userId)
		}

		if roleId := r.Header.Get(common.HeaderXRoleID); roleId != common.EmptyString {
			ctx = context.WithValue(ctx, common.CtxKeyRoleID, roleId)
		}

		// tenantctx.WithTenantIdCtx 内部同时写入 context.Value 与 outgoing metadata
		if tenantId := r.Header.Get(common.HeaderXTenantID); tenantId != common.EmptyString {
			ctx = tenantctx.WithTenantIdCtx(ctx, tenantId)
		}

		next(w, r.WithContext(ctx))
	}
}

// UserContext 返回可直接用于 rest.Server.Use 的中间件函数, 方便一行调用:
//
//	server.Use(middleware.UserContext())
func UserContext() rest.Middleware {
	return NewUserContextMiddleware().Handle
}
