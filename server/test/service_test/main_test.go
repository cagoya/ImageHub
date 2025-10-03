package test

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"server/service"
	"testing"
)

var db *gorm.DB
var userService *service.UserService
var imageService *service.ImageService
var dbConfig *DBConfig

func TestMain(m *testing.M) {
	// 1. 加载数据库配置
	var err error
	dbConfig, err = LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("dbConfig:", dbConfig)
	// 2. 创建测试数据库
	err = CreateDB(dbConfig)
	if err != nil {
		panic(err)
	}
	// 3. 连接测试数据库
	db, err = ConnectTestDB(dbConfig)
	if err != nil {
		panic(err)
	}
	// 4. 创建服务
	userService = service.NewUserService(db)
	imageService = service.NewImageService(db)
	// 5. 运行所有测试用例，并获取退出码
	code := m.Run()
	// 6. 删除测试数据库
	err = DropTestDB(dbConfig)
	if err != nil {
		panic(err)
	}
	// 7. 删除public目录
	err = os.RemoveAll("./public")
	if err != nil {
		panic(err)
	}
	// 8. 退出
	os.Exit(code)
}
