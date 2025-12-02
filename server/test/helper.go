package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"server/model"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/qdrant/go-client/qdrant"
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

type QdrantConfig struct {
	Host string
	Port int
}

// 从配置文件加载配置
func LoadConfig() (*DBConfig, *RedisConfig, *QdrantConfig, error) {
	dbConfig := &DBConfig{}
	redisConfig := &RedisConfig{}
	qdrantConfig := &QdrantConfig{}
	viper.SetConfigName("test.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, nil, err
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
	qdrantConfig.Host = viper.GetString("qdrant.host")
	qdrantConfig.Port = viper.GetInt("qdrant.port")
	return dbConfig, redisConfig, qdrantConfig, nil
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

// 连接测试 Qdrant 并创建 images_test collection
func ConnectTestQdrant(qdrantConfig *QdrantConfig) (*qdrant.Client, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantConfig.Host,
		Port: qdrantConfig.Port,
	})
	if err != nil {
		return nil, fmt.Errorf("连接 Qdrant 失败：%v", err)
	}

	// 检查 collection 是否存在
	ctx := context.Background()
	collectionName := "images_test"
	exist, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("检查集合是否存在失败：%v", err)
	}

	// 如果 collection 已存在，先删除它
	if exist {
		err = client.DeleteCollection(ctx, collectionName)
		if err != nil {
			return nil, fmt.Errorf("删除已存在的集合失败：%v", err)
		}
	}

	// 创建新的 collection
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1152,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("创建集合失败：%v", err)
	}

	return client, nil
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
