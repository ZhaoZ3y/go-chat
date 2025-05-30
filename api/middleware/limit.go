package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 用户级别的限流器
var userLimiters = make(map[int64]*rate.Limiter)

// GlobalRateLimit 全局限流中间件
func GlobalRateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), 100) // 每秒100个请求

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserRateLimit 用户级别限流中间件
func UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("User-ID")
		if userIDStr == "" {
			c.Next()
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.Next()
			return
		}

		limiter, exists := userLimiters[userID]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 每秒10个请求
			userLimiters[userID] = limiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "User rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
