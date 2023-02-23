package gerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Unknow            = status.New(codes.Unknown, "服务器异常").Err()
	ErrUnauthorized   = newError(10000, "请重新登录")
	ErrBadRequest     = newError(10001, "错误的请求")
	ErrBadCode        = newError(10002, "错误的验证码")
	ErrInGroup        = newError(10003, "不在聊天群组中")
	ErrGroupNotExist  = newError(10004, "群组不存在")
	ErrDeviceNotExist = newError(10005, "设备不存在")
	ErrAreadyIsFriend = newError(10006, "对方已经是好友")
	ErrUserNotFound   = newError(10007, "用户不存在")
	ErrNotFound       = newError(10404, "结果不存在")
)

func newError(code int, msg string) error {
	return status.New(codes.Code(code), msg).Err()
}
