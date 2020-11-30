package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const CtxUserIdKey = "userId"

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUser 获取当前用户登录的用户ID
func getCurrentUserId(ctx *gin.Context) (userId int64, err error) {
	uid, ok := ctx.Get(CtxUserIdKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userId, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return

}
