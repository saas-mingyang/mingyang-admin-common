package auth

import (
	"context"
	"errors"
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
	Authorization = "authorization"
)

type Claims struct {
	UserId   string `json:"userId"`
	RoleId   string `json:"roleId"`
	DeptId   int64  `json:"deptId"`
	TenantId int64  `json:"jwtTenantId"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

// UnaryAuthInterceptor 返回 RPC 拦截器
func UnaryAuthInterceptor(skipMethods []string, secretKey string) grpc.UnaryServerInterceptor {
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
			logx.Errorf("failed to get metadata from context")
			return nil, ErrUnauthorized
		}
		tokens := md.Get(Authorization)
		if len(tokens) <= 0 {
			return nil, ErrUnauthorized
		}
		_, err := ValidateToken(tokens[0], secretKey)
		if err != nil {
			logx.Errorf("unauthorized err = %s\n", err)
			return nil, ErrUnauthorized
		}
		return handler(ctx, req)
	}
}

func ValidateToken(token, secretKey string) (Claims, error) {
	var claims Claims
	if token == "" {
		return claims, errors.New("token is empty")
	}
	fromToken := jwt.StripBearerPrefixFromToken(token)
	if fromToken == "" {
		return claims, errors.New("token is empty")
	}
	jwtToken, err := jwt.ParseJwtToken(fromToken, secretKey)
	if err != nil {
		return claims, errors.New("ParseJwtToken error")
	}
	err = jwt.MapClaimsToStruct(jwtToken, &claims)
	if err != nil {
		logx.Errorf("MapClaimsToStruct error = %v", err)
		return claims, errors.New("MapClaimsToStruct error")
	}
	return claims, nil
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
	ctx := r.Context()
	token := r.Header.Get(Authorization)

	// 追加 token 到 metadata（不丢失原有 metadata）
	ctxWithTokenMD := metadata.AppendToOutgoingContext(ctx, Authorization, token)

	// 将 token 也放入 context.Value 以便当前进程使用
	return context.WithValue(ctxWithTokenMD, Authorization, token)
}
