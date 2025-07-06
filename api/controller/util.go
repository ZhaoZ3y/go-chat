package controller

import (
	"IM/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

// getAndParseUserID 是一个辅助函数，用于从Gin上下文中获取并解析用户ID
func getAndParseUserID(c *gin.Context) (int64, bool) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return 0, false
	}
	userId, ok := userIDAny.(int64)
	if !ok {
		response.ClientErrorResponse(c, response.ParamErrorCode, "用户ID类型错误")
		return 0, false
	}
	return userId, true
}
