package helper

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host string
	Port string
	Password string
	Db int
}

// 清空Redis中的内容
func ClearDB(RedisConfig *RedisConfig) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisConfig.Host, RedisConfig.Port),
		Password: RedisConfig.Password,
		DB:       RedisConfig.Db,
	})
	_, err := redisClient.FlushAll(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("清空Redis失败：%v", err)
	}
	return nil
}