package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/utils/jwt"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

var ErrUnauthorized = status.Error(codes.Unauthenticated, "unauthorized")

// Authorization AUTH_JWT_TOKEN 常量定义
const (
	Authorization = "Authorization" // g
)

// UnaryAuthInterceptor 返回 RPC 拦截器
func UnaryAuthInterceptor(skipMethods []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if containsMethod(skipMethods, info.FullMethod) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			fmt.Printf("获取md失败")
			return nil, ErrUnauthorized
		}
		tokens := md.Get(Authorization)
		if len(tokens) <= 0 {
			fmt.Printf("当前方法: %s\n", info.FullMethod)
			fmt.Println("tokens len 0")
			return nil, ErrUnauthorized
		}
		// 调用 AuthService 验证 token
		_, err := ValidateToken(tokens[0])
		if err != nil {
			logx.Errorf("unauthorized err = %s\n", err)
			return nil, ErrUnauthorized
		}
		return handler(ctx, req)
	}
}

func ValidateToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("token is empty")
	}
	fromToken := jwt.StripBearerPrefixFromToken(token)
	fmt.Printf("fromToken = %s\n", fromToken)
	return "", nil
}

func containsMethod(methods []string, target string) bool {
	for _, m := range methods {
		if m == target {
			return true
		}
	}
	return false
}

// SetTokenToContext 将token设置到Context,作为RPC认证的方式
func SetTokenToContext(r *http.Request) context.Context {
	token := r.Header.Get(Authorization)
	return metadata.NewOutgoingContext(r.Context(), metadata.Pairs(Authorization, token))
}
