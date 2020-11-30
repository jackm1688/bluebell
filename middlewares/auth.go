package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//1客户端携带token有三种方式:1.放在请求头 2.放在请求体 3.放在URI
		//这里假设Token放在Header的Authorization中，并使用Bearer开头
		//Authorization: Bearer xxx.xx.xxx
		//这里具体使用实现方式要依赖于实际的业务情况决定
		authorHeadr := ctx.Request.Header.Get("Authorization")
		if authorHeadr == "" {
			controller.ResponseError(ctx, controller.CodeNeedLogin)
			ctx.Abort()
			return
		}
		//按空格分割
		parts := strings.SplitN(authorHeadr, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(ctx, controller.CodeInvalidToken)
			ctx.Abort()
			return
		}
		// parts[1]是获取到的tokenString,我们使用之前定义好的解析JWT函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(ctx, controller.CodeInvalidToken)
			ctx.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		ctx.Set(controller.CtxUserIdKey, mc.UserId)
		ctx.Next() // 后续的处理函数可以用过c.Get("userId")来获取当前请求的用户信息
	}
}
