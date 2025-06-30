package controller

import (
	"IM/api/rpc"
	"IM/pkg/model/request"
	"IM/pkg/utils/response"
	"IM/rpc/user/user"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.UserClient.Register(ctx, &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	})

	if err != nil {
		logx.Errorf("注册失败: %v", err)
		response.ServerErrorResponse(c, "注册失败: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.UserClient.Login(ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		logx.Errorf("登录失败: %v", err)
		response.ClientErrorResponse(c, response.RPCClientErrorCode, "登录失败: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		response.ClientErrorResponse(c, response.ParamErrorCode, "用户ID不能为空")
		return
	}
	var userIdInt int64
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "用户ID格式错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.UserClient.GetUserInfo(ctx, &user.GetUserInfoRequest{
		UserId: userIdInt,
	})

	if err != nil {
		logx.Errorf("获取用户信息失败: %v", err)
		response.ServerErrorResponse(c, "获取用户信息失败: "+err.Error())
		return
	}

	response.SuccessResponse(c, resp.UserInfo)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}
	userId, ok := userIDAny.(int64)
	if !ok {
		response.ClientErrorResponse(c, response.ParamErrorCode, "用户ID类型错误")
		return
	}

	var req request.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.UserClient.UpdateUserInfo(ctx, &user.UpdateUserInfoRequest{
		UserId:   userId,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Email:    req.Email,
	})

	if err != nil {
		logx.Errorf("更新用户信息失败: %v", err)
		response.ServerErrorResponse(c, "更新用户信息失败: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.ClientErrorResponse(c, response.UnauthorizedCode, "用户未登录")
		return
	}
	userId, ok := userIDAny.(int64)
	if !ok {
		response.ClientErrorResponse(c, response.ParamErrorCode, "用户ID类型错误")
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.UserClient.ChangePassword(ctx, &user.ChangePasswordRequest{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		logx.Errorf("修改密码失败: %v", err)
		response.ServerErrorResponse(c, "修改密码失败: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": resp.Message,
	})
}

// RefreshToken 刷新Token
func RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.UserClient.RefreshToken(ctx, &user.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		logx.Errorf("刷新Token失败: %v", err)
		response.ServerErrorResponse(c, "刷新Token失败: "+err.Error())
		return
	}

	response.SuccessResponse(c, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}
