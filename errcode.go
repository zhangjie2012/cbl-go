package cbl

var (
	ErrNone                = ""                      // 没错
	ErrBadRequest          = "Bad Request"           // 错误的请求（参数错误）
	ErrLoginRequired       = "Login Required"        // 需要登录
	ErrPermissionDenied    = "Permission Denied"     // 权限不足
	ErrInternalServerError = "Internal Server Error" // 服务器内部错误, 服务器自身逻辑问题
	ErrNotFound            = "Not Found"             // 资源不存在
	ErrOutOfRange          = "Out Of Range"          // 越界访问
	ErrTooManyRequests     = "Too Many Requests"     // 请求过于频繁, 限流
)
