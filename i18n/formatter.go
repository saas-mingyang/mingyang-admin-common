package i18n

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//////////////////// Formatter ////////////////////

// Formatter 国际化消息格式化器
type Formatter struct {
	translator Translator
}

// NewFormatter 创建 Formatter
func NewFormatter(translator Translator) *Formatter {
	return &Formatter{translator: translator}
}

// Format 普通国际化格式化（支持任意参数）
// 示例：formatter.Format(ctx, "user.not_found", userId)
func (f *Formatter) Format(ctx context.Context, key string, args ...any) string {
	msg := f.translator.Trans(ctx, key)
	if len(args) == 0 {
		return msg
	}
	return formatArgs(msg, args...)
}

// FormatNamed 命名参数格式化
// 示例：formatter.FormatNamed(ctx, "user.not_found", map[string]any{"name":"张三"})
func (f *Formatter) FormatNamed(ctx context.Context, key string, params map[string]any) string {
	msg := f.translator.Trans(ctx, key)
	return formatNamedArgs(msg, params)
}

// NewError 创建国际化错误
func (f *Formatter) NewError(ctx context.Context, key string, args ...any) *I18nError {
	msg := f.Format(ctx, key, args...)
	return &I18nError{
		Key:     key,
		Message: msg,
	}
}

// NewGrpcError 创建 gRPC 国际化错误（强烈推荐）
func (f *Formatter) NewGrpcError(ctx context.Context, code codes.Code, key string, args ...any) error {
	msg := f.Format(ctx, key, args...)
	return status.Error(code, msg)
}

// FromGrpcError 从 gRPC 错误提取消息
func (f *Formatter) FromGrpcError(err error) string {
	if st, ok := status.FromError(err); ok {
		return st.Message()
	}
	return err.Error()
}

//////////////////// I18nError ////////////////////

// I18nError 国际化错误结构
type I18nError struct {
	Key     string
	Message string
}

func (e *I18nError) Error() string {
	return e.Message
}

//////////////////// 内部工具函数 ////////////////////

// formatArgs 数字占位符格式化：{0} {1}...
func formatArgs(message string, args ...any) string {
	for i, arg := range args {
		ph := "{" + strconv.Itoa(i) + "}"
		message = strings.ReplaceAll(message, ph, fmt.Sprint(arg))
	}
	return message
}

// formatNamedArgs 命名占位符格式化：{name} {age}...
func formatNamedArgs(message string, params map[string]any) string {
	for k, v := range params {
		message = strings.ReplaceAll(message, "{"+k+"}", fmt.Sprint(v))
	}
	return message
}
