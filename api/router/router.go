package router

import (
	"IM/api/controller"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/register", controller.Register)             // 用户注册
		user.POST("/login", controller.Login)                   // 用户登录
		user.GET("/info/:id", controller.GetUserInfo)           // 获取用户信息
		user.PUT("/update", controller.UpdateUserInfo)          // 更新用户信息
		user.PUT("/update/password", controller.ChangePassword) // 更新用户密码
		user.GET("/search", controller.SearchUser)              // 搜索用户
		user.POST("/refresh_token", controller.RefreshToken)    // 刷新用户令牌
	}

	return r
}
