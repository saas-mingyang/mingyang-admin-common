package valid

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/nyaruka/phonenumbers"
)

type ContactType int

const (
	Unknown ContactType = iota
	Mobile
	Email
)

// 编译一次，多次使用
var (
	emailRegex         *regexp.Regexp
	phoneRegexCache    sync.Map // 缓存电话号码解析结果
	once               sync.Once
	commonCountryCodes = []string{"US", "CN", "GB", "JP", "KR", "DE", "FR", "IN", "BR", "RU"}
)

// EmailValidator 邮箱验证器
type EmailValidator struct {
	allowUnicode   bool
	validateMX     bool // 是否验证MX记录
	allowIPAddress bool // 是否允许IP地址作为域名
}

// NewEmailValidator 创建邮箱验证器
func NewEmailValidator(options ...func(*EmailValidator)) *EmailValidator {
	v := &EmailValidator{
		allowUnicode:   false,
		validateMX:     false,
		allowIPAddress: false,
	}
	for _, option := range options {
		option(v)
	}
	return v
}

// WithUnicode 允许Unicode字符
func WithUnicode(allow bool) func(*EmailValidator) {
	return func(v *EmailValidator) {
		v.allowUnicode = allow
	}
}

// WithMXValidation 启用MX记录验证
func WithMXValidation(enable bool) func(*EmailValidator) {
	return func(v *EmailValidator) {
		v.validateMX = enable
	}
}

// IsValidEmail 验证邮箱格式（标准库方式）
func IsValidEmail(email string) bool {
	return IsValidEmailWithOptions(email)
}

// IsValidEmailWithOptions 可配置的邮箱验证
func IsValidEmailWithOptions(email string, options ...func(*EmailValidator)) bool {
	if email == "" || !utf8.ValidString(email) {
		return false
	}

	// 基本长度检查
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// 使用标准库验证
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// 使用正则验证格式
	if !emailRegex.MatchString(email) {
		return false
	}

	// 分离本地部分和域名
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]

	// 检查本地部分长度
	if len(local) > 64 {
		return false
	}

	// 检查域名长度
	if len(domain) > 255 {
		return false
	}

	// 创建验证器
	validator := NewEmailValidator(options...)

	// 验证Unicode
	if !validator.allowUnicode && containsUnicode(email) {
		return false
	}

	// 验证IP地址域名
	if !validator.allowIPAddress && isIPAddress(domain) {
		return false
	}

	return true
}

// PhoneValidator 电话号码验证器
type PhoneValidator struct {
	defaultRegion   string
	allowedTypes    []phonenumbers.PhoneNumberType
	lenientParsing  bool
	cacheEnabled    bool
	validationCache *sync.Map
}

type cacheEntry struct {
	valid   bool
	number  *phonenumbers.PhoneNumber
	country string
}

// NewPhoneValidator 创建电话号码验证器
func NewPhoneValidator(defaultRegion string, options ...func(*PhoneValidator)) *PhoneValidator {
	v := &PhoneValidator{
		defaultRegion: strings.ToUpper(defaultRegion),
		allowedTypes: []phonenumbers.PhoneNumberType{
			phonenumbers.MOBILE,
			phonenumbers.FIXED_LINE_OR_MOBILE,
		},
		lenientParsing:  true,
		cacheEnabled:    true,
		validationCache: &sync.Map{},
	}
	for _, option := range options {
		option(v)
	}
	return v
}

// WithAllowedTypes 设置允许的电话类型
func WithAllowedTypes(types []phonenumbers.PhoneNumberType) func(*PhoneValidator) {
	return func(v *PhoneValidator) {
		v.allowedTypes = types
	}
}

// WithLenientParsing 设置宽松解析模式
func WithLenientParsing(lenient bool) func(*PhoneValidator) {
	return func(v *PhoneValidator) {
		v.lenientParsing = lenient
	}
}

// WithCache 启用/禁用缓存
func WithCache(enabled bool) func(*PhoneValidator) {
	return func(v *PhoneValidator) {
		v.cacheEnabled = enabled
	}
}

// IsValidPhone 验证电话号码
func IsValidPhone(phone, countryCode string) bool {
	validator := NewPhoneValidator(countryCode)
	return validator.Validate(phone)
}

// Validate 验证电话号码
func (v *PhoneValidator) Validate(phone string) bool {
	if phone == "" {
		return false
	}

	phone = strings.TrimSpace(phone)

	// 检查缓存
	if v.cacheEnabled {
		if cached, ok := v.validationCache.Load(phone); ok {
			if entry, ok := cached.(cacheEntry); ok && entry.country == v.defaultRegion {
				return entry.valid
			}
		}
	}

	// 尝试多种方式解析
	parsed, err := v.parsePhoneNumber(phone)
	if err != nil {
		if v.cacheEnabled {
			v.validationCache.Store(phone, cacheEntry{valid: false, country: v.defaultRegion})
		}
		return false
	}

	// 验证号码
	valid := phonenumbers.IsValidNumber(parsed)
	if !valid {
		if v.cacheEnabled {
			v.validationCache.Store(phone, cacheEntry{valid: false, country: v.defaultRegion})
		}
		return false
	}

	// 检查号码类型
	numType := phonenumbers.GetNumberType(parsed)
	typeAllowed := false
	for _, allowedType := range v.allowedTypes {
		if numType == allowedType {
			typeAllowed = true
			break
		}
	}

	if !typeAllowed {
		if v.cacheEnabled {
			v.validationCache.Store(phone, cacheEntry{valid: false, country: v.defaultRegion})
		}
		return false
	}

	// 缓存结果
	if v.cacheEnabled {
		v.validationCache.Store(phone, cacheEntry{
			valid:   true,
			number:  parsed,
			country: v.defaultRegion,
		})
	}

	return true
}

// parsePhoneNumber 解析电话号码
func (v *PhoneValidator) parsePhoneNumber(phone string) (*phonenumbers.PhoneNumber, error) {
	var parsed *phonenumbers.PhoneNumber
	var err error

	// 策略1：如果包含+号，尝试国际格式解析
	if strings.HasPrefix(phone, "+") {
		parsed, err = phonenumbers.Parse(phone, "")
		if err == nil {
			return parsed, nil
		}
	}

	// 策略2：使用指定的国家代码
	if v.defaultRegion != "" {
		parsed, err = phonenumbers.Parse(phone, v.defaultRegion)
		if err == nil {
			return parsed, nil
		}
	}

	// 策略3：尝试常见国家代码
	for _, country := range commonCountryCodes {
		parsed, err = phonenumbers.Parse(phone, country)
		if err == nil && phonenumbers.IsValidNumber(parsed) {
			return parsed, nil
		}
	}

	return nil, fmt.Errorf("unable to parse phone number: %s", phone)
}

// GetPhoneInfo 获取电话号码详细信息
func GetPhoneInfo(phone, countryCode string) (*PhoneInfo, error) {
	validator := NewPhoneValidator(countryCode)
	return validator.GetPhoneInfo(phone)
}

func (v *PhoneValidator) GetPhoneInfo(phone string) (*PhoneInfo, error) {
	parsed, err := v.parsePhoneNumber(phone)
	if err != nil {
		return nil, err
	}

	regionCode := phonenumbers.GetRegionCodeForNumber(parsed)
	countryCode := parsed.GetCountryCode()
	numType := phonenumbers.GetNumberType(parsed)

	info := &PhoneInfo{
		InternationalFormat: phonenumbers.Format(parsed, phonenumbers.INTERNATIONAL),
		NationalFormat:      phonenumbers.Format(parsed, phonenumbers.NATIONAL),
		E164Format:          phonenumbers.Format(parsed, phonenumbers.E164),
		CountryCode:         countryCode,
		RegionCode:          regionCode,
		IsValid:             phonenumbers.IsValidNumber(parsed),
		IsMobile:            numType == phonenumbers.MOBILE || numType == phonenumbers.FIXED_LINE_OR_MOBILE,
		Type:                int(phonenumbers.GetNumberType(parsed)),
		NationalNumber:      phonenumbers.GetNationalSignificantNumber(parsed),
	}

	return info, nil
}

// PhoneInfo 电话号码信息
type PhoneInfo struct {
	InternationalFormat string
	NationalFormat      string
	E164Format          string
	CountryCode         int32
	RegionCode          string
	IsValid             bool
	IsMobile            bool
	Type                int
	NationalNumber      string
}

// CheckContactType 判断输入是邮箱还是手机号
func CheckContactType(input, countryCode string) ContactType {
	return CheckContactTypeWithOptions(input, countryCode)
}

// CheckContactTypeWithOptions 可配置的检查联系方式类型
func CheckContactTypeWithOptions(input, countryCode string, options ...func(*PhoneValidator)) ContactType {
	input = strings.TrimSpace(input)
	if input == "" {
		return Unknown
	}

	// 先检查是否为邮箱（速度更快）
	if IsValidEmail(input) {
		return Email
	}

	// 检查是否为手机号
	validator := NewPhoneValidator(countryCode, options...)
	if validator.Validate(input) {
		return Mobile
	}

	return Unknown
}

// CleanPhoneNumber 清理电话号码
func CleanPhoneNumber(phone string) string {
	// 移除所有非数字和+字符
	var builder strings.Builder
	builder.Grow(len(phone))

	for _, r := range phone {
		if (r >= '0' && r <= '9') || r == '+' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// 辅助函数
func containsUnicode(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}

func isIPAddress(s string) bool {
	// 简单的IP地址检查
	if strings.Count(s, ".") == 3 {
		parts := strings.Split(s, ".")
		if len(parts) == 4 {
			for _, part := range parts {
				if part == "" || len(part) > 3 {
					return false
				}
				for _, r := range part {
					if r < '0' || r > '9' {
						return false
					}
				}
			}
			return true
		}
	}
	return false
}

// BatchValidatePhones 批量验证电话号码
func BatchValidatePhones(phones []string, countryCode string) map[string]bool {
	results := make(map[string]bool)
	validator := NewPhoneValidator(countryCode, WithCache(true))

	for _, phone := range phones {
		results[phone] = validator.Validate(phone)
	}

	return results
}

// ValidateAndNormalizePhone 验证并标准化电话号码
func ValidateAndNormalizePhone(phone, countryCode string) (normalized string, isValid bool) {
	validator := NewPhoneValidator(countryCode)

	if !validator.Validate(phone) {
		return "", false
	}

	info, err := validator.GetPhoneInfo(phone)
	if err != nil {
		return "", false
	}

	return info.E164Format, true
}
