package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//然后我们定义JWT的过期时间，这里以2小时为例：
var TokenExpreDuration = time.Hour * 2

//接下来还需要定义Secret：
var MySecrect = []byte("夏天夏天悄悄过去")

//生成JWT
// GenToken 生成JWT
func GetToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := &MyClaims{
		Username: "username",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpreDuration).Unix(),
			Issuer:    "my-project",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodES256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecrect)
}

//解析JWT
// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecrect, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); !ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
