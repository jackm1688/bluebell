package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

//过期时间
//const TokenExpireDuration = time.Hour * 2

//加掩
var mySecret = []byte("夏天夏天悄悄过去")

// GenToken 生成JWT
func GenToken(userId int64, usernmae string) (string, error) {
	//创建一个我们自己的明文
	c := MyClaims{
		UserId:   userId,
		Username: usernmae, //自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), //过期时间
			Issuer: "bluebell", //签发人
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用指定的secret签名获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	//校验token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
