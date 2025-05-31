package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SuccessCode        = 20000 // 成功响应码
	ParamErrorCode     = 40000 // 参数错误响应码
	ServerErrorCode    = 50000 // 服务器错误响应码
	RPCClientErrorCode = 50001 // RPC客户端错误响应码
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
