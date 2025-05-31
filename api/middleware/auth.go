package middleware

import (
	"IM/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 Authorization 头部
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少或非法的 token"})
			c.Abort()
			return
		}

		// 提取 token 字符串
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 token
		claims, err := jwt.ParseAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 token"})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
