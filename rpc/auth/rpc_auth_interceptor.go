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

type Claims struct {
	UserId   int64  `json:"userId"`
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
			fmt.Printf("获取md失败")
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
	fmt.Printf("fromToken = %s\n", fromToken)
	if fromToken == "" {
		return claims, errors.New("token is empty")
	}
	jwtToken, err := jwt.ParseJwtToken(fromToken, secretKey)
	if err != nil {
		return claims, errors.New("ParseJwtToken error")
	}
	fmt.Printf("jwtToken = %v\n", jwtToken)
	err = jwt.MapClaimsToStruct(jwtToken, &claims)
	if err != nil {
		fmt.Printf("MapClaimsToStruct error = %v\n", err)
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
	token := r.Header.Get(Authorization)
	return metadata.NewOutgoingContext(r.Context(), metadata.Pairs(Authorization, token))
}
