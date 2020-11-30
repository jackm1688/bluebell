package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

//存放业务逻辑

// Signup 用户注册
func Signup(p *models.ParamSigup) (err error) {
	//1.判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//2.生成uid
	userId := snowflake.GetID()
	u := models.User{
		UserId:   userId,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存进数据库
	return mysql.InsertUser(&u)
}

// Login 用户登录
func Login(u *models.User) (user *models.User, err error) {
	user = &models.User{
		UserId:   u.UserId,
		Username: u.Username,
	}
	//1.直接登录
	if err := mysql.Login(u); err != nil {
		return nil, err
	}
	//生成jwt的token
	token, err := jwt.GenToken(u.UserId, u.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
