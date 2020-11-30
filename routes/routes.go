package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "bluebell/docs"

	"github.com/gin-contrib/pprof"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Setup(cfg *settings.AppConfig) *gin.Engine {
	if cfg.App.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	r.Use(logger.GinLogger(), logger.GinReconvery(true))
	r.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	r.Static("/static", "./static")
	r.LoadHTMLFiles("templates/index.html")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	//注册业务流程
	v1 := r.Group("/api/v1")
	v1.Use(middlewares.Cors())
	{
		v1.POST("/signup", controller.SignupHandler)
		v1.POST("/login", controller.LoginHandler)
	}

	//根据时间或分数进行获取帖子
	v1.GET("/posts2", controller.GetPostListHandler2)

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		//获取评论列表
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDeailHandler)
		v1.POST("/post", controller.CratePostHandler)
		v1.GET("/post/:id", controller.PostDetailHandler)

		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	pprof.Register(r) //注册pprof
	return r
}
