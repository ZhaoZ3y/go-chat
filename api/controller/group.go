package controller

import (
	"IM/api/rpc"
	"IM/pkg/model/request"
	"IM/pkg/utils/response"
	"IM/rpc/group/group"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateGroup 创建群组
func CreateGroup(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.CreateGroup(c.Request.Context(), &group.CreateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		OwnerId:     userId.(int64),
		MemberIds:   req.MemberIds,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// GetGroupInfo 获取群组信息
func GetGroupInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	groupIdStr := c.Param("group_id")
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "群组ID格式错误")
		return
	}

	resp, err := rpc.GroupClient.GetGroupInfo(c.Request.Context(), &group.GetGroupInfoRequest{
		GroupId: groupId,
		UserId:  userId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if resp.GroupInfo == nil {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, "群组不存在")
		return
	}

	response.SuccessResponse(c, gin.H{
		"group_info":       resp.GroupInfo,
		"user_member_info": resp.UserMemberInfo,
	})
}

// GetGroupList 获取群组列表
func GetGroupList(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	resp, err := rpc.GroupClient.GetGroupList(c.Request.Context(), &group.GetGroupListRequest{
		UserId: userId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"groups": resp.Groups,
		"total":  resp.Total,
	})
}

// GetGroupMemberList 获取群组成员列表
func GetGroupMemberList(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
	}

	groupIdStr := c.Param("group_id")
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "群组ID格式错误")
		return
	}

	resp, err := rpc.GroupClient.GetGroupMemberList(c.Request.Context(), &group.GetGroupMemberListRequest{
		GroupId: groupId,
		UserId:  userId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if len(resp.Members) == 0 {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, "群组不存在或无成员")
		return
	}

	response.SuccessResponse(c, gin.H{
		"members": resp.Members,
		"total":   resp.Total,
	})

}

// UpdateGroupInfo 更新群组信息
func UpdateGroupInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.UpdateGroupInfo(c.Request.Context(), &group.UpdateGroupInfoRequest{
		GroupId:     req.GroupId,
		OperatorId:  userId.(int64),
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// SetMemberRole 设置群组成员角色
func SetMemberRole(c *gin.Context) {
	operatorId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
	}

	var req request.SetMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.SetMemberRole(c.Request.Context(), &group.SetMemberRoleRequest{
		GroupId:    req.GroupId,
		OperatorId: operatorId.(int64),
		UserId:     req.UserId,
		Role:       req.Role,
	})
	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// MuteMember 禁言群成员
func MuteMember(c *gin.Context) {
	operatorId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.MuteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.MuteMember(c.Request.Context(), &group.MuteMemberRequest{
		GroupId:    req.GroupId,
		OperatorId: operatorId.(int64),
		UserId:     req.UserId,
		Duration:   req.Duration,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// InviteToGroup 邀请进群
func InviteToGroup(c *gin.Context) {
	inviterId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.InviteToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.InviteToGroup(c.Request.Context(), &group.InviteToGroupRequest{
		GroupId:   req.GroupId,
		InviterId: inviterId.(int64),
		UserIds:   req.UserIds,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message":         resp.Message,
		"failed_user_ids": resp.FailedUserIds,
	})
}

// JoinGroup 加入群组
func JoinGroup(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.JoinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.JoinGroup(c.Request.Context(), &group.JoinGroupRequest{
		GroupId: req.GroupId,
		UserId:  userId.(int64),
		Reason:  req.Reason,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// KickFromGroup 踢出群成员
func KickFromGroup(c *gin.Context) {
	operatorId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.KickFromGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.KickFromGroup(c.Request.Context(), &group.KickFromGroupRequest{
		GroupId:    req.GroupId,
		OperatorId: operatorId.(int64),
		UserId:     req.UserId,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// LeaveGroup 退出群组
func LeaveGroup(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.LeaveGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.LeaveGroup(c.Request.Context(), &group.LeaveGroupRequest{
		GroupId: req.GroupId,
		UserId:  userId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// DismissGroup 解散群组
func DismissGroup(c *gin.Context) {
	ownerId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.DismissGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.DismissGroup(c.Request.Context(), &group.DismissGroupRequest{
		GroupId: req.GroupId,
		OwnerId: ownerId.(int64),
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// TransferGroup 转让群组
func TransferGroup(c *gin.Context) {
	ownerId, exists := c.Get("user_id")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "未授权访问")
		return
	}

	var req request.TransferGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	resp, err := rpc.GroupClient.TransferGroup(c.Request.Context(), &group.TransferGroupRequest{
		GroupId:    req.GroupId,
		OwnerId:    ownerId.(int64),
		NewOwnerId: req.NewOwnerId,
	})

	if err != nil {
		response.ServerErrorResponse(c, "服务器内部错误: "+err.Error())
		return
	}

	if !resp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, resp.Message)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}
