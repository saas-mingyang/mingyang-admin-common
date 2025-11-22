package auth

import (
	"errors"
	"fmt"
)

type RpcAuthService struct {
	// 可以加入 Redis/DB/JWT 等依赖
}

// NewAuthService 初始化
func NewAuthService() *RpcAuthService {
	return &RpcAuthService{}
}

// ValidateToken 验证 token 返回 userID
func (a *RpcAuthService) ValidateToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("token empty")
	}
	fmt.Printf("token = %s\n", token)

	return token, nil
}
