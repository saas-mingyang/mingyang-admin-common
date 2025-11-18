package str

import (
	"strings"
	"unicode"
)

// IsNotEmpty 检查字符串是否不为空
func IsNotEmpty(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Concat 拼接多个字符串
func Concat(strs ...string) string {
	var builder strings.Builder
	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}

// Join 使用指定分隔符连接字符串切片
func Join(separator string, strs ...string) string {
	return strings.Join(strs, separator)
}

// Contains 检查字符串是否包含子串（忽略大小写）
func Contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// StartsWith 检查字符串是否以指定前缀开始
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith 检查字符串是否以指定后缀结束
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// ToUpper 转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower 转换为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Trim 去除字符串首尾空白字符
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimPrefix 去除指定前缀
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// TrimSuffix 去除指定后缀
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// Replace 替换字符串中的子串
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Split 分割字符串
func Split(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

// Substring 截取子串
func Substring(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return ""
	}
	return s[start:end]
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// RemoveSpaces 移除所有空格
func RemoveSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

// CamelToSnake 驼峰命名转蛇形命名
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// SnakeToCamel 蛇形命名转驼峰命名
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
