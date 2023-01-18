package app

import (
	"context"
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"github.com/go-redis/redis/v9"
)

func NewCache() *redis.Client {
	appConfig := config.NewAppConfig()
	client := redis.NewClient(&redis.Options{
		Password: appConfig.RedisPassword,
		Addr:     appConfig.RedisHost + ":" + appConfig.RedisPort,
	})
	var ctx = context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis connected")
	return client

}
