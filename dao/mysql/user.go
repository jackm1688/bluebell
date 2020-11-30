package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"
)

var secret = "*&^%mfzLe^%$mon@"

//把每一步封装成函数
//待Logic层调用

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) error {

	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	fmt.Println("COUNT:", count)
	if count > 0 {
		return ErrUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(u *models.User) (err error) {
	//生成密码
	u.Password = encryptPassword(u.Password)
	//执行sql语句入库
	sqlStr := "insert into user(user_id,username,password) VALUES(?,?,?)"
	_, err = db.Exec(sqlStr, u.UserId, u.Username, u.Password)
	return
}

func encryptPassword(pass string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(pass)))
}

// Login 用户登录校验
func Login(u *models.User) error {
	//生成密码
	oPass := u.Password
	sqlStr := "select user_id,username,password from user where username=?"
	if err := db.Get(u, sqlStr, u.Username); err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotExist
		}
		return err
	}

	if u.Password != encryptPassword(oPass) {
		return ErrUserPasswdWrong
	}
	return nil
}

func GetUserById(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id,username from user where user_id=?"
	zap.L().Debug("", zap.String("SQL", sqlStr), zap.Any("uid", id))
	err = db.Get(user, sqlStr, id)
	return

}
