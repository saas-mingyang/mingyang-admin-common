package common

// HTTP 状态码定义

const (
	// 1xx Informational
	StatusContinue           int = 100
	StatusSwitchingProtocols int = 101
	StatusProcessing         int = 102
	StatusEarlyHints         int = 103

	// 2xx Success
	StatusOK                   int = 200
	StatusCreated              int = 201
	StatusAccepted             int = 202
	StatusNonAuthoritativeInfo int = 203
	StatusNoContent            int = 204
	StatusResetContent         int = 205
	StatusPartialContent       int = 206
	StatusMultiStatus          int = 207
	StatusAlreadyReported      int = 208
	StatusIMUsed               int = 226

	// 3xx Redirection
	StatusMultipleChoices   int = 300
	StatusMovedPermanently  int = 301
	StatusFound             int = 302
	StatusSeeOther          int = 303
	StatusNotModified       int = 304
	StatusUseProxy          int = 305
	StatusTemporaryRedirect int = 307
	StatusPermanentRedirect int = 308

	// 4xx Client Error
	StatusBadRequest                   int = 400
	StatusUnauthorized                 int = 401
	StatusPaymentRequired              int = 402
	StatusForbidden                    int = 403
	StatusNotFound                     int = 404
	StatusMethodNotAllowed             int = 405
	StatusNotAcceptable                int = 406
	StatusProxyAuthRequired            int = 407
	StatusRequestTimeout               int = 408
	StatusConflict                     int = 409
	StatusGone                         int = 410
	StatusLengthRequired               int = 411
	StatusPreconditionFailed           int = 412
	StatusRequestEntityTooLarge        int = 413
	StatusRequestURITooLong            int = 414
	StatusUnsupportedMediaType         int = 415
	StatusRequestedRangeNotSatisfiable int = 416
	StatusExpectationFailed            int = 417
	StatusTeapot                       int = 418
	StatusMisdirectedRequest           int = 421
	StatusUnprocessableEntity          int = 422
	StatusLocked                       int = 423
	StatusFailedDependency             int = 424
	StatusTooEarly                     int = 425
	StatusUpgradeRequired              int = 426
	StatusPreconditionRequired         int = 428
	StatusTooManyRequests              int = 429
	StatusRequestHeaderFieldsTooLarge  int = 431
	StatusUnavailableForLegalReasons   int = 451

	// 5xx Server Error
	StatusInternalServerError           int = 500
	StatusNotImplemented                int = 501
	StatusBadGateway                    int = 502
	StatusServiceUnavailable            int = 503
	StatusGatewayTimeout                int = 504
	StatusHTTPVersionNotSupported       int = 505
	StatusVariantAlsoNegotiates         int = 506
	StatusInsufficientStorage           int = 507
	StatusLoopDetected                  int = 508
	StatusNotExtended                   int = 510
	StatusNetworkAuthenticationRequired int = 511
)
