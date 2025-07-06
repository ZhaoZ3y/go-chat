package controller

import (
	"IM/api/rpc"
	"IM/pkg/model/request"
	"IM/pkg/utils/response"
	"IM/rpc/friend/friend"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// SendFriendRequest 发送好友申请
func SendFriendRequest(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	var req request.SendFriendRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.SendFriendRequest(rpcCtx, &friend.SendFriendRequestRequest{
		FromUserId: userID,
		ToUserId:   req.ToUserID,
		Message:    req.Message,
	})
	if err != nil {
		logx.Errorf("RPC SendFriendRequest failed: %v", err)
		response.ServerErrorResponse(c, "发送好友申请失败")
		return
	}
	if !rpcResp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, rpcResp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message":    "好友申请已发送",
		"request_id": rpcResp.RequestId,
	})
}

// HandleFriendRequest 处理好友申请
func HandleFriendRequest(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	var req request.HandleFriendRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.HandleFriendRequest(rpcCtx, &friend.HandleFriendRequestRequest{
		RequestId: req.RequestID,
		UserId:    userID,
		Action:    req.Action,
		Message:   req.Message,
	})
	if err != nil {
		logx.Errorf("RPC HandleFriendRequest failed: %v", err)
		response.ServerErrorResponse(c, "处理好友申请失败")
		return
	}
	if !rpcResp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, rpcResp.Message)
		return
	}

	// 异步任务触发点：
	// 在一个完整的系统中，可以在这里通过 goroutine 触发一个异步任务，
	// 例如使用WebSocket通知申请方，申请已被处理。
	// go notifyFriendRequestResult(rpcResp.RequestInfo)

	response.SuccessResponse(c, gin.H{
		"message": rpcResp.Message,
	})
}

// GetFriendList 获取当前用户的好友列表
func GetFriendList(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.GetFriendList(rpcCtx, &friend.GetFriendListRequest{
		UserId: userID,
	})
	if err != nil {
		logx.Errorf("RPC GetFriendList failed: %v", err)
		response.ServerErrorResponse(c, "获取好友列表失败")
		return
	}

	response.SuccessResponse(c, rpcResp)
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	var req request.FriendDeleteActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.DeleteFriend(rpcCtx, &friend.DeleteFriendRequest{
		UserId:   userID,
		FriendId: req.FriendID,
	})
	if err != nil {
		logx.Errorf("RPC DeleteFriend failed: %v", err)
		response.ServerErrorResponse(c, "删除好友失败")
		return
	}
	if !rpcResp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, rpcResp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "删除好友成功",
	})
}

// BlockFriend 拉黑或取消拉黑好友 (toggle)
func BlockFriend(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	var req request.BlockFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.BlockFriend(rpcCtx, &friend.BlockFriendRequest{
		UserId:   userID,
		FriendId: req.FriendID,
		Status:   req.Status,
	})
	if err != nil {
		logx.Errorf("RPC BlockFriend failed: %v", err)
		response.ServerErrorResponse(c, "操作失败")
		return
	}
	if !rpcResp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, rpcResp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": rpcResp.Message, // 消息会是 "拉黑好友成功" 或 "已取消拉黑"
	})
}

// UpdateFriendRemark 更新好友备注
func UpdateFriendRemark(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	var req request.UpdateFriendRemarkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.UpdateFriendRemark(rpcCtx, &friend.UpdateFriendRemarkRequest{
		UserId:   userID,
		FriendId: req.FriendID,
		Remark:   req.Remark,
	})
	if err != nil {
		logx.Errorf("RPC UpdateFriendRemark failed: %v", err)
		response.ServerErrorResponse(c, "更新备注失败")
		return
	}
	if !rpcResp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, rpcResp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "更新备注成功",
	})
}

// GetBlockedFriendList 获取黑名单列表
func GetBlockedFriendList(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FriendClient.GetBlockedList(rpcCtx, &friend.GetBlockedListRequest{
		UserId: userID,
	})
	if err != nil {
		logx.Errorf("RPC GetBlockedList failed: %v", err)
		response.ServerErrorResponse(c, "获取黑名单列表失败")
		return
	}

	response.SuccessResponse(c, rpcResp)
}
