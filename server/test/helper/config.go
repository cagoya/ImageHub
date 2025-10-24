package helper

import (
	"github.com/spf13/viper"
)

// 从配置文件加载配置
func LoadConfig() (*DBConfig, *RedisConfig, error) {
	dbConfig := &DBConfig{}
	redisConfig := &RedisConfig{}
	viper.SetConfigName("test.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, err
	}
	dbConfig.Username = viper.GetString("db.username")
	dbConfig.Password = viper.GetString("db.password")
	dbConfig.Host = viper.GetString("db.host")
	dbConfig.Port = viper.GetString("db.port")
	dbConfig.Database = viper.GetString("db.database")
	redisConfig.Host = viper.GetString("redis.host")
	redisConfig.Port = viper.GetString("redis.port")
	redisConfig.Password = viper.GetString("redis.password")
	redisConfig.Db = viper.GetInt("redis.db")
	return dbConfig, redisConfig, nil
}