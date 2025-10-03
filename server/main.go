package main

import (
	"os"
	"path/filepath"

	"server/api"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 读取配置文件
	err := global.ReadConfigFile()
	if err != nil {
		panic(err)
	}
	// 初始化zap
	err = global.InitZap()
	if err != nil {
		panic(err)
	}
	logger := global.GetLogger()
	defer logger.Sync() // 确保程序退出时同步日志
	logger.Info("zap初始化成功")

	// 初始化db
	err = global.InitDB()
	if err != nil {
		panic(err)
	}
	logger.Info("db初始化成功")
	// 初始化redis
	global.InitRedis()
	logger.Info("redis初始化成功")

	// 创建 Gin 路由
	r := gin.Default()

	// 添加日志中间件
	r.Use(middleware.LoggerMiddleware())

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 添加静态文件服务
	// 创建静态文件目录
	dir := filepath.Dir("./public")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	r.Static("/public", "./public")

	apiGroup := r.Group("/api")
	userHandler := api.NewUserHandler()
	imageHandler := api.NewImageHandler()

	// 用户相关无需认证的接口
	userGroup := apiGroup.Group("/user")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}

	// 用户相关需要认证的接口
	userGroup.Use(middleware.JWTMiddleware())
	{
		userGroup.GET("/getUserById", userHandler.GetUserByID)
		userGroup.GET("/getUserByUsername", userHandler.GetUserByUsername)
		userGroup.GET("/getAllUsers", userHandler.GetAllUsers)
		userGroup.PUT("/updateUser", userHandler.UpdateUser)
		userGroup.DELETE("/deleteUser", userHandler.DeleteUser)
		userGroup.PUT("/grantAdminPermission", userHandler.GrantAdminPermission)
		userGroup.PUT("/revokeAdminPermission", userHandler.RevokeAdminPermission)
	}

	// 图片相关接口
	imageGroup := apiGroup.Group("/image").Use(middleware.JWTMiddleware())
	{
		imageGroup.POST("/upload", imageHandler.UploadImage)
	}

	// 启动服务
	port := viper.GetString("server.port")
	logger.Info("服务器启动", zap.String("port", port))
	r.Run("localhost:" + port)
}
