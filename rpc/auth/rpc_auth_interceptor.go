package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/datapermctx"
	"github.com/saas-mingyang/mingyang-admin-common/utils/jwt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/enum"
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
	XDataScope    = "x-data-scope"
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

// TransferHTTPContextToGRPC 将HTTP请求的上下文中的多个值传递到gRPC调用的metadata中
func TransferHTTPContextToGRPC(r *http.Request) context.Context {
	ctx := r.Context()

	// 准备metadata
	md := metadata.MD{}

	// 1. 传递认证token
	token := r.Header.Get(Authorization)
	if token != "" {
		md.Set(Authorization, token)
	}

	// 2. 传递数据权限标识
	// 先尝试从上下文中获取
	if dataScope, ok := ctx.Value(datapermctx.ScopeKey).(string); ok {
		md.Set(string(datapermctx.ScopeKey), dataScope)
	} else {
		// 如果上下文中没有，尝试从HTTP头部获取
		dataScope := r.Header.Get(XDataScope)
		if dataScope != "" {
			md.Set(string(datapermctx.ScopeKey), dataScope)
		}
	}

	// 3. 传递租户ID
	if tenantID, ok := ctx.Value(enum.TenantIdCtxKey).(string); ok && tenantID != "" {
		md.Set(enum.TenantIdCtxKey, tenantID)
	}

	// 4. 传递用户ID
	if userID, ok := ctx.Value(enum.UserIdRpcCtxKey).(string); ok && userID != "" {
		md.Set(enum.UserIdRpcCtxKey, userID)
	}

	// 5. 传递角色ID
	if roleID, ok := ctx.Value(enum.RoleIdRpcCtxKey).(string); ok && roleID != "" {
		md.Set(enum.RoleIdRpcCtxKey, roleID)
	}

	// 6. 传递部门ID
	if deptID, ok := ctx.Value(enum.DepartmentIdRpcCtxKey).(string); ok && deptID != "" {
		md.Set(enum.DepartmentIdRpcCtxKey, deptID)
	}

	// 7. 传递语言
	if lang, ok := ctx.Value(enum.I18nCtxKey).(string); ok && lang != "" {
		md.Set(enum.I18nCtxKey, lang)
	}

	// 8. 传递客户端IP
	if clientIP, ok := ctx.Value(enum.ClientIPCtxKey).(string); ok && clientIP != "" {
		md.Set(enum.ClientIPCtxKey, clientIP)
	}

	// 如果metadata不为空，则将其附加到上下文
	if len(md) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}
