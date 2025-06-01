package router

import (
	"IM/api/controller"
	"IM/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) // 跨域中间件

	r.POST("/register", controller.Register)          // 用户注册
	r.POST("/login", controller.Login)                // 用户登录
	r.POST("/refresh_token", controller.RefreshToken) // 刷新用户令牌

	// 用户相关路由
	user := r.Group("/user")
	{
		user.GET("/info", middleware.AuthMiddleware(), controller.GetUserInfo)               // 获取用户信息
		user.PUT("/update/info", middleware.AuthMiddleware(), controller.UpdateUserInfo)     // 更新用户信息
		user.PUT("/update/password", middleware.AuthMiddleware(), controller.ChangePassword) // 更新用户密码
	}

	r.GET("/search", controller.Search) // 搜索用户和群组

	// 群组相关路由
	group := r.Group("/group")
	{
		group.POST("/create", middleware.AuthMiddleware(), controller.CreateGroup)         // 创建群组
		group.GET("/info", middleware.AuthMiddleware(), controller.GetGroupInfo)           // 获取群组信息
		group.GET("/list", middleware.AuthMiddleware(), controller.GetGroupList)           // 获取用户的群组列表
		group.GET("/members", middleware.AuthMiddleware(), controller.GetGroupMemberList)  // 获取群组成员列表
		group.PUT("/update/info", middleware.AuthMiddleware(), controller.UpdateGroupInfo) // 更新群组信息
		group.POST("/set/role", middleware.AuthMiddleware(), controller.SetMemberRole)     // 设置群组角色
		group.POST("/mute/member", middleware.AuthMiddleware(), controller.MuteMember)     // 禁言群成员
		group.POST("/invite", middleware.AuthMiddleware(), controller.InviteToGroup)       // 邀请用户加入群组
		group.POST("/join", middleware.AuthMiddleware(), controller.JoinGroup)             // 用户加入群组
		group.POST("/leave", middleware.AuthMiddleware(), controller.LeaveGroup)           // 用户退出群组
		group.POST("/dismiss", middleware.AuthMiddleware(), controller.DismissGroup)       // 解散群组
		group.POST("/transfer", middleware.AuthMiddleware(), controller.TransferGroup)     // 转让群组
		group.POST("/kick", middleware.AuthMiddleware(), controller.KickFromGroup)         // 踢出群成员
	}
	return r
}
