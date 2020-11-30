package test

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/settings"
	"fmt"

	"go.uber.org/zap"
)

func Init() {
	//1.初始化配置文件
	if err := settings.Init("/Users/mengfanzhen/go/src/bluebell/conf/config.yaml"); err != nil {
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
}
