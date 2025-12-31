package valid

import (
	"regexp"
	"strings"
)

type ContactType int

const (
	Unknown ContactType = iota
	Mobile
	Email
)

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	// RFC 5322 标准邮箱正则表达式
	pattern := `(?:[a-z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// IsInternationalMobile IsValidPhone 验证手机号格式
func IsInternationalMobile(mobile string) bool {
	// 国际手机号：可选+号开头，后跟1-15位数字
	pattern := `^\+?[1-9]\d{1,14}$`
	matched, err := regexp.MatchString(pattern, mobile)
	if err != nil {
		return false
	}
	return matched
}

// CheckContactType 判断输入是邮箱还是手机号
func CheckContactType(input string) ContactType {
	input = strings.TrimSpace(input)

	// 先检查是否为手机号
	if IsInternationalMobile(input) || IsInternationalMobile(input) {
		return Mobile
	}

	// 再检查是否为邮箱
	if IsValidEmail(input) {
		return Email
	}

	return Unknown
}
