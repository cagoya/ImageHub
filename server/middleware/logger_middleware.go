package middleware

import (
	"time"

	"server/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware gin框架的路由日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	logger := global.GetLogger()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算请求耗时
		cost := time.Since(start)

		// 记录日志
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		}

		// 如果有错误，添加错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		// 根据状态码选择日志级别
		if c.Writer.Status() >= 400 {
			logger.Error("HTTP Request", fields...)
		} else {
			logger.Info("HTTP Request", fields...)
		}
	}
}
