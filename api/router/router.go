package router

import (
	"IM/api/controller"
	"IM/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controller.Register)          // 用户注册
	r.POST("/login", controller.Login)                // 用户登录
	r.POST("/refresh_token", controller.RefreshToken) // 刷新用户令牌
	// 用户相关路由
	user := r.Group("/user")
	{
		user.GET("/info", middleware.AuthMiddleware(), controller.GetUserInfo)               // 获取用户信息
		user.PUT("/update/info", middleware.AuthMiddleware(), controller.UpdateUserInfo)     // 更新用户信息
		user.PUT("/update/password", middleware.AuthMiddleware(), controller.ChangePassword) // 更新用户密码
		user.GET("/search", middleware.AuthMiddleware(), controller.SearchUser)              // 搜索用户
	}

	return r
}
