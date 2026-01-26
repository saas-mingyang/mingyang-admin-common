package i18n

import (
	"context"
	"strconv"
	"strings"
)

// Formatter 国际化消息格式化器
type Formatter struct {
	translator *Translator
}

// NewFormatter 创建新的格式化器
func NewFormatter(translator *Translator) *Formatter {
	return &Formatter{
		translator: translator,
	}
}

// FormatMessage 格式化国际化消息（兼容您提供的签名）
func (f *Formatter) FormatMessage(ctx context.Context,
	messageKey string, args ...string) string {

	// 获取原始消息
	message := f.translator.Trans(ctx, messageKey)

	// 如果没有参数，直接返回
	if len(args) == 0 {
		return message
	}

	// 替换占位符 {0}, {1}...
	for i, arg := range args {
		placeholder := "{" + strconv.Itoa(i) + "}"
		message = strings.ReplaceAll(message, placeholder, arg)
	}

	return message
}

// FormatMessageWithInterface 支持任意类型参数
func (f *Formatter) FormatMessageWithInterface(ctx context.Context,
	messageKey string, args ...interface{}) string {

	// 获取原始消息
	message := f.translator.Trans(ctx, messageKey)

	// 如果没有参数，直接返回
	if len(args) == 0 {
		return message
	}

	// 替换占位符 {0}, {1}...
	for i, arg := range args {
		placeholder := "{" + strconv.Itoa(i) + "}"
		argStr := toString(arg)
		message = strings.ReplaceAll(message, placeholder, argStr)
	}

	return message
}

// FormatError 格式化错误消息
func (f *Formatter) FormatError(translator *Translator, ctx context.Context,
	messageKey string, args ...string) error {

	message := f.FormatMessage(ctx, messageKey, args...)
	return NewI18nError(messageKey, message)
}

// FormatErrorWithInterface 格式化错误消息（支持任意类型）
func (f *Formatter) FormatErrorWithInterface(ctx context.Context,
	messageKey string, args ...interface{}) error {

	message := f.FormatMessageWithInterface(ctx, messageKey, args...)
	return NewI18nError(messageKey, message)
}

// toString 将任意类型转换为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(v.(int64), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return ""
	}
}

// I18nError 国际化错误
type I18nError struct {
	Key     string
	Message string
}

// NewI18nError 创建国际化错误
func NewI18nError(key, message string) *I18nError {
	return &I18nError{
		Key:     key,
		Message: message,
	}
}

func (e *I18nError) Error() string {
	return e.Message
}
