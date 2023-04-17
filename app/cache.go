package app

import (
	"context"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

func NewCache(appConfig *config.AppConfig, log *logrus.Logger) *redis.Client {
	cache := redis.NewClient(&redis.Options{
		Password: appConfig.RedisPassword,
		Addr:     appConfig.RedisHost + ":" + appConfig.RedisPort,
	})
	var ctx = context.Background()
	_, err := cache.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	cache.FlushAll(ctx)
	log.Info("redis connected successfully " + cache.String())
	return cache

}
