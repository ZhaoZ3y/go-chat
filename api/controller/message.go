package controller

import (
	"IM/api/rpc"
	"IM/pkg/model/request"
	"IM/pkg/utils/response"
	"IM/rpc/message/chat"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// SendMessage 发送信息
func SendMessage(c *gin.Context) {
	var req request.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	// 从 Gin 上下文中获取经过身份验证的用户ID
	fromUserID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.SendMessage(ctx, &chat.SendMessageRequest{
		FromUserId: fromUserID.(int64),
		ToUserId:   req.ToUserId,
		GroupId:    req.GroupId,
		Type:       chat.MessageType(req.Type),
		Content:    req.Content,
		Extra:      req.Extra,
		ChatType:   chat.ChatType(req.ChatType),
	})

	if err != nil {
		logx.Errorf("调用 SendMessage RPC 失败: %v", err)
		response.ServerErrorResponse(c, "发送消息失败")
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{"messageId": resp.MessageId})
}

// GetConversationList 获取会话列表
func GetConversationList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.GetConversationList(ctx, &chat.GetConversationListRequest{
		UserId: userID.(int64),
	})

	if err != nil {
		logx.Errorf("调用 GetConversationList RPC 失败: %v", err)
		response.ServerErrorResponse(c, "获取会话列表失败")
		return
	}

	response.SuccessResponse(c, resp.Conversations)
}

// GetMessageHistory 查看历史记录
func GetMessageHistory(c *gin.Context) {
	var req request.GetMessageHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.GetMessageHistory(ctx, &chat.GetMessageHistoryRequest{
		UserId:        userID.(int64),
		TargetId:      req.TargetId,
		ChatType:      chat.ChatType(req.ChatType),
		LastMessageId: req.LastMessageId,
		Limit:         req.Limit,
		Date:          req.Date,
	})

	if err != nil {
		logx.Errorf("调用 GetMessageHistory RPC 失败: %v", err)
		response.ServerErrorResponse(c, "获取消息历史失败")
		return
	}

	response.SuccessResponse(c, gin.H{
		"messages": resp.Messages,
		"hasMore":  resp.HasMore,
	})
}

// MarkMessageRead 标记已读
func MarkMessageRead(c *gin.Context) {
	var req request.MarkMessageReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.MarkMessageRead(ctx, &chat.MarkMessageReadRequest{
		UserId:            userID.(int64),
		TargetId:          req.TargetId,
		ChatType:          chat.ChatType(req.ChatType),
		LastReadMessageId: req.LastReadMessageId,
	})

	if err != nil {
		logx.Errorf("调用 MarkMessageRead RPC 失败: %v", err)
		response.ServerErrorResponse(c, "标记已读失败")
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, nil)
}

// RecallMessage 撤回消息
func RecallMessage(c *gin.Context) {
	var req request.RecallMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.RecallMessage(ctx, &chat.RecallMessageRequest{
		UserId:    userID.(int64),
		MessageId: req.MessageId,
	})

	if err != nil {
		logx.Errorf("调用 RecallMessage RPC 失败: %v", err)
		response.ServerErrorResponse(c, "撤回消息失败")
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, nil)
}

// DeleteMessage 删除消息
func DeleteMessage(c *gin.Context) {
	var req request.DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.DeleteMessage(ctx, &chat.DeleteMessageRequest{
		UserId:    userID.(int64),
		MessageId: req.MessageId,
	})

	if err != nil {
		logx.Errorf("调用 DeleteMessage RPC 失败: %v", err)
		response.ServerErrorResponse(c, "删除消息失败")
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, nil)
}

// DeleteConversation 删除会话
func DeleteConversation(c *gin.Context) {
	var req request.DeleteConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.DeleteConversation(ctx, &chat.DeleteConversationRequest{
		UserId:   userID.(int64),
		TargetId: req.TargetId,
		ChatType: chat.ChatType(req.ChatType),
	})

	if err != nil {
		logx.Errorf("调用 DeleteConversation RPC 失败: %v", err)
		response.ServerErrorResponse(c, "删除会话失败")
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, nil)
}

// PinConversation 置顶会话
func PinConversation(c *gin.Context) {
	var req request.PinConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.MessageClient.PinConversation(ctx, &chat.PinConversationRequest{
		UserId:   userID.(int64),
		TargetId: req.TargetId,
		ChatType: chat.ChatType(req.ChatType),
		IsPinned: req.IsPinned,
	})

	if err != nil {
		logx.Errorf("调用 PinConversation RPC 失败: %v", err)
		response.ServerErrorResponse(c, "操作失败")
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, nil)
}
