package test

import (
	"context"
	"fmt"
	"os"
	"server/global"
	"server/service"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB
var userService *service.UserService
var imageService *service.ImageService
var dbConfig *DBConfig
var redisConfig *RedisConfig
var qdrantConfig *QdrantConfig
var redisClient *redis.Client
var qdrantClient *qdrant.Client

func TestMain(m *testing.M) {
	// 1. 加载数据库、Redis 和 Qdrant 配置
	var err error
	dbConfig, redisConfig, qdrantConfig, err = LoadConfig()
	if err != nil {
		panic(err)
	}

	// 2. 连接测试数据库
	db, err = ConnectTestDB(dbConfig)
	if err != nil {
		panic(err)
	}

	// 3. 创建Redis客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})

	// 4. 测试Redis连接
	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("连接Redis失败: %v", err))
	}

	// 5. 清空Redis数据库
	err = ClearDB(redisConfig)
	if err != nil {
		panic(err)
	}

	// 5.5. 设置全局Redis客户端（用于向量化任务队列）
	global.SetRedisClient(redisClient)

	// 6. 连接 Qdrant 并创建 images_test collection
	qdrantClient, err = ConnectTestQdrant(qdrantConfig)
	if err != nil {
		panic(err)
	}

	// 7. 创建服务
	userService = service.NewUserService(db, redisClient)
	imageService = service.NewImageService(db, qdrantClient)

	// 8. 运行所有测试用例，并获取退出码
	code := m.Run()

	// 9. 清理资源并退出测试
	// 清空Redis数据库
	err = ClearDB(redisConfig)
	if err != nil {
		fmt.Printf("清空Redis失败: %v\n", err)
	}

	// 删除 Qdrant collection
	if qdrantClient != nil {
		ctx := context.Background()
		err = qdrantClient.DeleteCollection(ctx, "images_test")
		if err != nil {
			fmt.Printf("删除 Qdrant collection 失败: %v\n", err)
		}
	}

	// 关闭Redis连接
	redisClient.Close()

	// 删除public目录
	err = os.RemoveAll("./public")
	if err != nil {
		panic(err)
	}

	// 退出测试
	os.Exit(code)
}
