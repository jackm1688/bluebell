package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//处理注册请求的函数
func SignupHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var p = new(models.ParamSigup)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Signup with invalid param", zap.Error(err))
		//判断err是不是validator.ValidtionErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务处理
	if err := logic.Signup(p); err != nil {
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	//3.返回响应
	ResponseSucess(c, nil)
}

//LoginHandler 处理用户登录函数
func LoginHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var u = new(models.User)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator.ValidtionErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务处理
	user, err := logic.Login(u)
	if err != nil {
		zap.L().Warn("用户登录失败", zap.String("用户名", u.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrUserPasswdWrong) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	//3.返回响应
	ResponseSucess(c, gin.H{
		"user_id":  fmt.Sprintf("%d", user.UserId), //id值大于1 <<53-1  1 << 63-1
		"username": user.Username,
		"token":    user.Token,
	})
}
