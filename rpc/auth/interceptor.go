package auth

import (
	"context"
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/rpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var ErrUnauthorized = status.Error(codes.Unauthenticated, "unauthorized")

// UnaryAuthInterceptor 返回 RPC 拦截器
func UnaryAuthInterceptor(authSvc *RpcAuthService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 从 RPC Metadata 获取 token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrUnauthorized
		}
		tokens := md.Get("authorization")

		if len(tokens) <= 0 {
			fmt.Println("tokens len 0")
			return nil, ErrUnauthorized
		}
		// 调用 AuthService 验证 token
		userID, err := authSvc.ValidateToken(tokens[0])
		if err != nil {
			logx.Errorf("unauthorized err = %s\n", err)
			return nil, ErrUnauthorized
		}
		// 将 userID 写入 context，业务逻辑可直接使用
		ctx = context.WithValue(ctx, rpc.ContextKeyUserID, userID)
		return handler(ctx, req)
	}
}
