package oauthException

// 用于维护一些 `统一验证` 相关的 errCode
var (
	WrongAccount      = newError(501, "账号错误")
	WrongPassword     = newError(502, "密码错误")
	ClosedError       = newError(503, "统一系统在夜间关闭")
	NotActivatedError = newError(504, "账号未激活")
	OtherError        = newError(599, "其他错误")
)

// Error 自定义错误
type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

func newError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}
