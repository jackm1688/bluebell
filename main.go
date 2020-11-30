package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title blueblell
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path

func main() {

	var cfgFile = ""
	if len(os.Args) == 2 {
		cfgFile = os.Args[1]
	}

	//1.初始化配置文件
	if err := settings.Init(cfgFile); err != nil {
		panic(err)
	}
	fmt.Println("成功初始化配置文件")

	//2.初始化日志
	if err := logger.Init(settings.Conf); err != nil {
		panic(err)
	}
	fmt.Println("成功初始化日志实例")
	defer zap.L().Sync()

	//3.初始mysql实例
	if err := mysql.Init(settings.Conf); err != nil {
		panic(err)
	}
	defer mysql.Close()
	fmt.Println("初始化MySQL连接实例成功")

	//4.初始redis实例
	if err := redis.Init(settings.Conf); err != nil {
		panic(err)
	}
	defer redis.Close()
	fmt.Println("初始化Redis连接实例成功")

	//初始id生成器实例
	if err := snowflake.Init(settings.Conf.App.StartTime,
		settings.Conf.App.MachineId); err != nil {
		panic(err)
	}
	fmt.Println("初始化ID生成器实例成功")
	//初始化gin框架内校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		panic(err)
	}
	fmt.Println("初始化gin框架内校验器使用的翻译器成功")
	//5.初始化路由
	router := routes.Setup(settings.Conf)
	fmt.Println("初始化路由实例成功")
	//6.启动服务
	svc := http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.App.Port),
		Handler: router,
	}
	//开启一个goroutine
	go func() {
		zap.L().Debug("Start server")
		if err := svc.ListenAndServe(); err != nil {
			zap.L().Fatal("Listen: ", zap.Error(err))
		}
	}()
	fmt.Printf("成功启动[%s]:%d\n", settings.Conf.App.Name,
		settings.Conf.App.Port)

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quite := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quite, syscall.SIGINT, syscall.SIGTERM)
	<-quite // // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := svc.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
