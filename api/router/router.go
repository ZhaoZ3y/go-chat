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

	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("", hub.HandleWebSocket)
	}

	r.POST("/register", controller.Register)          // 用户注册
	r.POST("/login", controller.Login)                // 用户登录
	r.POST("/refresh_token", controller.RefreshToken) // 刷新用户令牌

	// 用户相关路由
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware()) // 需要身份验证的中间件
	{
		user.GET("/profile", controller.GetUserProfile)         // 获取用户个人资料
		user.GET("/info", controller.GetUserInfo)               // 获取用户信息
		user.PUT("/update/info", controller.UpdateUserInfo)     // 更新用户信息
		user.PUT("/update/password", controller.ChangePassword) // 更新用户密码
		user.PUT("/update/avatar", controller.UploadAvatar)     // 更新用户头像
	}

	r.GET("/search", middleware.AuthMiddleware(), controller.Search) // 搜索用户和群组

	// 好友相关路由
	friend := r.Group("/friend")
	friend.Use(middleware.AuthMiddleware())
	{
		friend.POST("/request/send", controller.SendFriendRequest)     // 发送好友申请
		friend.POST("/request/handle", controller.HandleFriendRequest) // 处理好友申请
		friend.GET("/list", controller.GetFriendList)                  // 获取好友列表
		friend.DELETE("/delete", controller.DeleteFriend)              // 删除好友
		friend.POST("/block", controller.BlockFriend)                  // 拉黑好友
		friend.GET("/blocked/list", controller.GetBlockedFriendList)   // 获取拉黑好友列表
		friend.PUT("/remark", controller.UpdateFriendRemark)           // 更新好友备注
	}

	// 群组相关路由
	group := r.Group("/group")
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/create", controller.CreateGroup)                        // 创建群组
		group.GET("/info", controller.GetGroupInfo)                          // 获取群组信息
		group.GET("/list", controller.GetGroupList)                          // 获取用户的群组列表
		group.GET("/members", controller.GetGroupMemberList)                 // 获取群组成员列表
		group.PUT("/update/info", controller.UpdateGroupInfo)                // 更新群组信息
		group.POST("/set/role", controller.SetMemberRole)                    // 设置群组角色
		group.POST("/mute/member", controller.MuteMember)                    // 禁言群成员
		group.POST("/invite", controller.InviteToGroup)                      // 邀请用户加入群组
		group.POST("/join", controller.JoinGroup)                            // 用户申请加入群组
		group.POST("/leave", controller.LeaveGroup)                          // 用户退出群组
		group.POST("/dismiss", controller.DismissGroup)                      // 解散群组
		group.POST("/transfer", controller.TransferGroup)                    // 转让群组
		group.POST("/kick", controller.KickFromGroup)                        // 踢出群成员
		group.GET("/member_info", controller.GetGroupMemberInfo)             // 获取群成员信息
		group.PUT("/update/member_info", controller.UpdateGroupMemberInfo)   // 更新群成员信息
		group.POST("/request/handle", controller.HandleJoinGroupApplication) // 处理入群申请
	}

	notification := r.Group("/notification")
	notification.Use(middleware.AuthMiddleware())
	{
		notification.GET("/request/friend", controller.GetFriendRequestList)                     // 获取好友申请列表
		notification.GET("/request/friend/unread-count", controller.GetUnreadFriendRequestCount) // 获取未读好友申请数量
		notification.GET("/group/notifications", controller.GetGroupNotifications)               // 获取群组相关通知
		notification.GET("/request/group", controller.GetJoinGroupApplications)                  // 获取入群申请列表
		notification.GET("/request/group/unread-count", controller.GetGroupUnreadCount)          // 获取未读入群申请数量
	}

	file := r.Group("/file")
	file.Use(middleware.AuthMiddleware())
	{
		file.POST("/upload", controller.UploadFile)    // 上传文件
		file.GET("/download", controller.DownloadFile) // 下载文件
		file.DELETE("/delete", controller.DeleteFile)  // 删除文件
		file.GET("/info", controller.GetFileInfo)      // 获取文件信息
		file.GET("/record", controller.GetFileRecord)  // 获取用户文件记录
	}

	chat := r.Group("/chat")
	chat.Use(middleware.AuthMiddleware()) // 需要身份验证的中间件
	{
		// 消息操作
		chat.POST("/send", controller.SendMessage)         // 发送消息 (POST /chat/send)
		chat.GET("/history", controller.GetMessageHistory) // 获取历史消息 (GET /chat/history)
		chat.POST("/read", controller.MarkMessageRead)     // 标记消息已读 (POST /chat/read)

		// 消息管理 (放在 /message 子路径下，更清晰)
		message := chat.Group("/message")
		{
			message.POST("/recall", controller.RecallMessage) // 撤回消息 (POST /chat/message/recall)
			message.POST("/delete", controller.DeleteMessage) // 删除消息 (POST /chat/message/delete)
		}

		// 会话操作 (放在 /conversation 子路径下)
		conversation := chat.Group("/conversation")
		{
			conversation.GET("/list", controller.GetConversationList)   // 获取会话列表 (GET /chat/conversation/list)
			conversation.POST("/delete", controller.DeleteConversation) // 删除会话 (POST /chat/conversation/delete)
			conversation.POST("/pin", controller.PinConversation)       // 置顶/取消置顶会话 (POST /chat/conversation/pin)
		}
	}

	return r
}
