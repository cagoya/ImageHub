package global

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

func InitRedis() error {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("连接Redis失败: %v", err)
	}
	return nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}
