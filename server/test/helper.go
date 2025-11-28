package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"server/model"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Db       int
}

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

// 清空数据库中的表
func ClearTables(dbConfig *DBConfig, tables []string) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}
	defer db.Close()

	// 测试连接
	if err = db.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}

	for _, table := range tables {
		_, err = db.Exec("DELETE FROM " + table)
		if err != nil {
			return fmt.Errorf("清空表失败：%v", err)
		}
	}
	return nil
}

// 连接测试数据库
func ConnectTestDB(dbConfig *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.User{}, &model.Image{}, &model.Tag{})
	if err != nil {
		return nil, fmt.Errorf("数据库迁移失败：%v", err)
	}

	return db, nil
}

// 清空Redis中的内容
func ClearDB(redisConfig *RedisConfig) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
	_, err := redisClient.FlushAll(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("清空Redis失败：%v", err)
	}
	return nil
}

// 复制文件
func CopyFile(src string, dst string, fileName string) error {
	// 读取源文件
	srcFile := path.Join(src, fileName)
	data, err := os.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("读取文件失败：%v", err)
	}

	// 写入目标文件
	dstFile := path.Join(dst, fileName)
	err = os.WriteFile(dstFile, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败：%v", err)
	}
	return nil
}
