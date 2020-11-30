package settings

import(
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//全局变量，保存程序所需的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	App   *App         `mapstructure:"app"`
	Log   *LogConfig   `mapstructure:"log"`
	MySQL *MySQLConfig `mapstructure:"mysql"`
	Redis *RedisConfig `mapstructure:"redis"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
	StartTime string `mapstructure:"star_time"`
	MachineId int64 `mapstructure:"machine_id"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge int `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	DriverName   string `mapstructure:"driver_name"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxConns     int    `mapstructure:"max_conn"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filename string) (err error)  {

	if filename != ""{
		viper.SetConfigFile(filename)
	}else{
		viper.SetConfigName("config")
		viper.AddConfigPath("./")
		viper.AddConfigPath("./conf")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("read config failed,err:%v\n",err)
		return err
	}

	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal config failed,err:%v\n",err)
		return err
	}

	//实时监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件更新了:%s\n",in.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal config failed,err:%v\n",err)
		}
	})
	return
}
