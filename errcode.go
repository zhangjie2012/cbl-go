package cbl

var (
	ErrNone                = ""                      // no error 没错
	ErrBadRequest          = "Bad Request"           // bad request 错误的请求（参数错误）
	ErrLoginRequired       = "Login Required"        // need login 需要登录
	ErrPermissionDenied    = "Permission Denied"     // permission denied 权限不足
	ErrInternalServerError = "Internal Server Error" // server logic error 服务器内部错误, 服务器自身逻辑问题
	ErrNotFound            = "Not Found"             // resource not found 资源不存在
	ErrOutOfRange          = "Out Of Range"          // out of range 越界访问
	ErrTooManyRequests     = "Too Many Requests"     // too many requests 请求过于频繁, 限流
)
