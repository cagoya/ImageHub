package main

import (
	"os"

	"server/api"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
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
	err = global.InitRedis()
	if err != nil {
		panic(err)
	}
	logger.Info("redis初始化成功")

	// 创建 Gin 路由
	r := gin.Default()

	// 添加日志中间件
	r.Use(middleware.LoggerMiddleware())

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 添加静态文件服务
	// 创建静态文件目录
	if err := os.MkdirAll("./public", 0755); err != nil {
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
		userGroup.POST("/forget-password", userHandler.ForgetPassword)
	}

	// 用户相关需要认证的接口
	userGroup.Use(middleware.JWTMiddleware())
	{
		userGroup.GET("/id", userHandler.GetUserByID)
		userGroup.GET("/username", userHandler.GetUserByUsername)
		userGroup.GET("/all", userHandler.GetAllUsers)
		userGroup.PUT("/update", userHandler.UpdateUser)
		userGroup.POST("/suspend", userHandler.SuspendUser)
		userGroup.POST("/unsuspend", userHandler.UnsuspendUser)
		userGroup.PUT("/grant", userHandler.GrantAdminPermission)
		userGroup.PUT("/revoke", userHandler.RevokeAdminPermission)
	}

	// 图片相关接口
	imageGroup := apiGroup.Group("/image").Use(middleware.JWTMiddleware())
	{
		imageGroup.POST("/upload", imageHandler.UploadImage)
		imageGroup.PUT("/update", imageHandler.UpdateImage)
		imageGroup.DELETE("/delete", imageHandler.DeleteImage)
		imageGroup.GET("/detail", imageHandler.GetImageDetail)
		imageGroup.GET("/file", imageHandler.GetImageFile)
		imageGroup.GET("/all", imageHandler.GetAllImages)
		imageGroup.GET("/user-id", imageHandler.GetImagesByUserID)
		imageGroup.GET("/tag", imageHandler.GetImagesByTag)
		imageGroup.GET("/time-range", imageHandler.GetImagesByTimeRange)
		imageGroup.GET("/area", imageHandler.GetImagesByArea)
		imageGroup.PUT("/tags", imageHandler.SetImageTags)
		imageGroup.GET("/generate-tags", imageHandler.GenerateImageTags)
	}

	// 启动服务
	logger.Info("服务器启动")
	r.Run("0.0.0.0:8888")
}
