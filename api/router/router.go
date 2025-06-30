package router

import (
	"IM/api/controller"
	"IM/api/middleware"
	"IM/pkg/websocket"
	"github.com/gin-gonic/gin"
)

func SetRouter(hub *websocket.Hub) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) // 跨域中间件

	r.POST("/register", controller.Register)          // 用户注册
	r.POST("/login", controller.Login)                // 用户登录
	r.POST("/refresh_token", controller.RefreshToken) // 刷新用户令牌

	// 用户相关路由
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware()) // 需要身份验证的中间件
	{
		user.GET("/info", controller.GetUserInfo)               // 获取用户信息
		user.PUT("/update/info", controller.UpdateUserInfo)     // 更新用户信息
		user.PUT("/update/password", controller.ChangePassword) // 更新用户密码
	}

	r.GET("/search", controller.Search) // 搜索用户和群组

	// 好友相关路由
	friend := r.Group("/friend")
	friend.Use(middleware.AuthMiddleware())
	{
		friend.POST("/request/send", controller.SendFriendRequest)                  // 发送好友申请
		friend.POST("/request/handle", controller.HandleFriendRequest)              // 处理好友申请
		friend.GET("/request/list", controller.GetFriendRequestList)                // 获取好友申请列表
		friend.GET("/request/unread-count", controller.GetUnreadFriendRequestCount) // 获取未读好友申请数量
		friend.GET("/list", controller.GetFriendList)                               // 获取好友列表
		friend.DELETE("/delete", controller.DeleteFriend)                           // 删除好友
		friend.POST("/block", controller.BlockFriend)                               // 拉黑好友
		friend.GET("/blocked/list", controller.GetBlockedFriendList)                // 获取拉黑好友列表
		friend.PUT("/remark", controller.UpdateFriendRemark)                        // 更新好友备注
	}

	// 群组相关路由
	group := r.Group("/group")
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/create", controller.CreateGroup)         // 创建群组
		group.GET("/info", controller.GetGroupInfo)           // 获取群组信息
		group.GET("/list", controller.GetGroupList)           // 获取用户的群组列表
		group.GET("/members", controller.GetGroupMemberList)  // 获取群组成员列表
		group.PUT("/update/info", controller.UpdateGroupInfo) // 更新群组信息
		group.POST("/set/role", controller.SetMemberRole)     // 设置群组角色
		group.POST("/mute/member", controller.MuteMember)     // 禁言群成员
		group.POST("/invite", controller.InviteToGroup)       // 邀请用户加入群组
		group.POST("/join", controller.JoinGroup)             // 用户加入群组
		group.POST("/leave", controller.LeaveGroup)           // 用户退出群组
		group.POST("/dismiss", controller.DismissGroup)       // 解散群组
		group.POST("/transfer", controller.TransferGroup)     // 转让群组
		group.POST("/kick", controller.KickFromGroup)         // 踢出群成员

	}
	return r
}
