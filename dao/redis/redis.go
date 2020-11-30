package redis

import (
	"bluebell/settings"
	"fmt"

	"github.com/go-redis/redis"
)

var client *redis.Client

func Init(cfg *settings.AppConfig) (err error) {

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if err = client.Ping().Err(); err != nil {
		return err
	}
	return
}

func Close() {
	_ = client.Close()
}
