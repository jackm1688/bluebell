package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(ctx *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(ctx *gin.Context) {
		//如果获取不到令牌就返回响应
		if bucket.TakeAvailable(1) == 0 {
			ctx.String(http.StatusOK, "rate limit...")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
