package middleware

import (
	"strings"

	"server/utils"

	"github.com/gin-gonic/gin"
)

type Claim struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

// JWT认证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 提取Bearer令牌
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			utils.Unauthorized(c, "无效的认证令牌格式")
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}
