package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"niuNiuSDKBackend/common/log"

	conf "niuNiuSDKBackend/config"
	"time"
)

var Rdb *redis.Client

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.GetConfig().RedisConfig.Addr,
		Password: conf.GetConfig().RedisConfig.Password, // 密码
		DB:       conf.GetConfig().RedisConfig.DB,       // 数据库
		PoolSize: conf.GetConfig().RedisConfig.PoolSize, // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Logger.Fatal("redis init err", log.Any("redis init err", err.Error()))
	}

}
