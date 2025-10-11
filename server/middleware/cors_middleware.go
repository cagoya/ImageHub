package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源，这里为了调试方便，暂时允许所有源
		c.Header("Access-Control-Allow-Origin", "*")

		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-auth-token")

		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

		// 设置暴露给前端的响应头（重要：允许前端读取自定义响应头）
		c.Header("Access-Control-Expose-Headers", "x-auth-token")

		// 设置允许的凭据
		c.Header("Access-Control-Allow-Credentials", "true")

		// 设置预检请求的缓存时间
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
