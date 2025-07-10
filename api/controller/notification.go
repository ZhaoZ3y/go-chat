package controller

import (
	"IM/api/rpc"
	"IM/pkg/model/request"
	"IM/pkg/utils/response"
	"IM/rpc/friend/friend"
	"IM/rpc/group/group"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// GetFriendRequestList 获取收到的好友申请列表
func GetFriendRequestList(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	var req request.GetFriendRequestListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.GetFriendRequestList(rpcCtx, &friend.GetFriendRequestListRequest{
		UserId: userID,
		Status: req.Status,
	})
	if err != nil {
		logx.Errorf("RPC GetFriendRequestList failed: %v", err)
		response.ServerErrorResponse(c, "获取申请列表失败")
		return
	}

	response.SuccessResponse(c, rpcResp)
}

// GetUnreadFriendRequestCount 获取未读好友申请数量
func GetUnreadFriendRequestCount(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.GetUnreadFriendRequestCount(rpcCtx, &friend.GetUnreadFriendRequestCountRequest{
		UserId: userID,
	})
	if err != nil {
		logx.Errorf("RPC GetUnreadFriendRequestCount failed: %v", err)
		response.ServerErrorResponse(c, "获取未读数量失败")
		return
	}

	response.SuccessResponse(c, gin.H{
		"count": rpcResp.Count,
	})
}

// GetGroupNotifications 获取用户的群组相关通知
func GetGroupNotifications(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	resp, err := rpc.GroupClient.GetGroupNotifications(c.Request.Context(), &group.GetGroupNotificationsRequest{
		UserId: userId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"notifications": resp.Notifications,
	})
}

// GetJoinGroupApplications 获取入群申请列表 (管理员或群主)
func GetJoinGroupApplications(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	// 权限验证在RPC层处理
	resp, err := rpc.GroupClient.GetJoinGroupApplications(c.Request.Context(), &group.GetJoinGroupApplicationsRequest{})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"applications": resp.Applications,
	})
}

// GetGroupUnreadCount 获取未读群组通知总数
func GetGroupUnreadCount(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	resp, err := rpc.GroupClient.GetUnreadCount(c.Request.Context(), &group.GetUnreadCountRequest{
		UserId: userId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"total_unread_count": resp.TotalUnreadCount,
	})
}
