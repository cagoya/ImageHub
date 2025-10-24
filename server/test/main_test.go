package test

import (
	"context"
	"fmt"
	"os"
	"server/service"
	"server/test/helper"
	"testing"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB
var userService *service.UserService
var imageService *service.ImageService
var dbConfig *helper.DBConfig
var redisConfig *helper.RedisConfig
var redisClient *redis.Client

func TestMain(m *testing.M) {
	// 1. 加载数据库和Redis配置
	var err error
	dbConfig, redisConfig, err = helper.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("dbConfig:", dbConfig)
	fmt.Println("redisConfig:", redisConfig)

	// 2. 创建测试数据库
	err = helper.CreateDB(dbConfig)
	if err != nil {
		panic(err)
	}

	// 3. 连接测试数据库
	db, err = helper.ConnectTestDB(dbConfig)
	if err != nil {
		panic(err)
	}

	// 4. 创建真实Redis客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})

	// 5. 测试Redis连接
	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("连接Redis失败: %v", err))
	}

	// 6. 清空Redis数据库
	err = helper.ClearDB(redisConfig)
	if err != nil {
		panic(err)
	}

	// 7. 创建服务
	userService = service.NewUserService(db, redisClient)
	imageService = service.NewImageService(db)

	// 8. 运行所有测试用例，并获取退出码
	code := m.Run()

	// 9. 清理资源
	// 清空Redis数据库
	err = helper.ClearDB(redisConfig)
	if err != nil {
		fmt.Printf("清空Redis失败: %v\n", err)
	}

	// 关闭Redis连接
	redisClient.Close()

	// 删除测试数据库
	err = helper.DropTestDB(dbConfig)
	if err != nil {
		panic(err)
	}

	// 删除public目录
	err = os.RemoveAll("./public")
	if err != nil {
		panic(err)
	}

	// 10. 退出
	os.Exit(code)
}
