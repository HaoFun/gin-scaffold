package bootstrap

import (
	"context"
	"gin-scaffold/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitializeRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     global.App.Config.Redis.Host + ":" + global.App.Config.Redis.Port,
		Password: global.App.Config.Redis.Password,
		DB:       global.App.Config.Redis.DB,
		PoolSize: global.App.Config.Redis.PoolSize,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.App.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
		return nil
	}
	return client
}
