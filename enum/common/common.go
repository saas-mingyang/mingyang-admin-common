package common

const (
	StatusNormal uint8 = 1
	StatusBanned uint8 = 2
)

const (
	DefaultParentId uint64 = 1000000
	// DefaultInvalidRoleId used to judge whether the token belongs to core, if it is DefaultInvalidRoleId, it belongs to mms
	DefaultInvalidRoleId uint64 = 1000000

	// DateFormat 时间格式
	DateFormat      = "2006-01-02"
	DateTimeFormat  = "2006-01-02 15:04:05"
	TimeFormat      = "15:04:05"
	TimeZone        = "Asia/Shanghai"
	TimeStampFormat = "20060102150405"

	// 符号

	// Comma 逗号
	Comma = ","
	// Semicolon 分号
	Semicolon = ";"
	// Space 空格
	Space = " "
	// VerticalLine 竖线
	VerticalLine = "|"
	// Colon 冒号
	Colon = ":"
	// Dot 点
	Dot = "."
	// Slash 斜杠
	Slash = "/"
	// BackSlash 反斜杠
	BackSlash = "\\"
	// At 符号
	At = "@"
	// Underscore 下划线
	Underscore = "_"
	// Minus 减号
	Minus = "-"
	// Plus 加号
	Plus = "+"
	// Star 星号
	Star = "*"
	// QuestionMark ？
	QuestionMark = "?"
	// Hash 井号
	Hash = "#"
	// Percent 符号
	Percent = "%"
	// Ampersand 符号
	Ampersand = "&"
	// Pipe 符号
	Pipe = "|"
	// Equal 等于
	Equal = "="
	// EmptyString ""
	EmptyString = ""

	// ContentTypeJson application/json
	ContentTypeJson = "application/json"
	// ContentTypeForm application/x-www-form-urlencoded
	ContentTypeForm = "application/x-www-form-urlencoded"
	// ContentTypeMultipartForm multipart/form-data
	ContentTypeMultipartForm = "multipart/form-data"
	// ContentTypeTextPlain text/plain
	ContentTypeTextPlain = "text/plain"
	// ContentTypeTextHtml text/html
	ContentTypeTextHtml = "text/html"
	// ContentTypeTextXml text/xml
	ContentTypeTextXml = "text/xml"
	// ContentType ContentTypeOctetStream application/octet-stream
	ContentType = "Content-Type"
	//XForwardedFor  X-Forwarded-For
	XForwardedFor = "X-Forwarded-For"
	// XRealIP X-Real-IP
	XRealIP = "X-Real-IP"
	// BearerPrefix "Bearer "
	BearerPrefix = "Bearer"
	// AcceptLanguage Accept-Language"
	AcceptLanguage = "Accept-Language"
	// ContentLength Content-Length"
	ContentLength = "Content-Length"
	// ContentDisposition Content-Disposition
	ContentDisposition     = "Content-Disposition"
	AcceptRanges           = "Accept-Ranges"
	ContentTypeOctetStream = "application/octet"

	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"

	Zero    = 0
	One     = 1
	Two     = 2
	Three   = 3
	Four    = 4
	Five    = 5
	Six     = 6
	Seven   = 7
	Eight   = 8
	Nine    = 9
	Ten     = 10
	Hundred = 100
	KB      = 1024
)
