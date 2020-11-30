package mysql

import "errors"

var (
	ErrUserExist       = errors.New("用户已经存在")
	ErrUserNotExist    = errors.New("用户已不存在")
	ErrUserPasswdWrong = errors.New("密码不正确")
	ErrorInvalidID     = errors.New("无效的ID")
)
