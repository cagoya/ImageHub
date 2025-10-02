package main

import (
	"server/global"

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
	// 创建 Gin 路由
	r := gin.Default()

	// 添加静态文件服务
	r.Static("/public", "./public")

	// apiGroup := r.Group("/api")

	// 启动服务
	port := viper.GetString("server.port")
	logger.Info("服务器启动", zap.String("port", port))
	r.Run("127.0.0.1:" + port)
}
