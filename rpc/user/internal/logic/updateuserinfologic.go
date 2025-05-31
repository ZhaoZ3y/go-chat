package logic

import (
	"IM/pkg/model"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoRequest) (*user.UpdateUserInfoResponse, error) {
	if in.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID无效")
	}

	// 检查用户是否存在
	var existUser model.User
	err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", in.UserId).First(&existUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 如果更新邮箱，检查邮箱是否已被其他用户使用
	if in.Email != "" && in.Email != existUser.Email {
		var emailUser model.User
		err = l.svcCtx.DB.Where("email = ? AND id != ? AND deleted_at IS NULL", in.Email, in.UserId).First(&emailUser).Error
		if err == nil {
			return nil, status.Error(codes.AlreadyExists, "邮箱已被使用")
		}
		if err != gorm.ErrRecordNotFound {
			l.Logger.Errorf("查询邮箱失败: %v", err)
			return nil, status.Error(codes.Internal, "系统错误")
		}
	}

	// 构建更新数据
	updateData := map[string]interface{}{
		"update_at": time.Now().Unix(),
	}

	if in.Nickname != "" {
		updateData["nickname"] = in.Nickname
	}
	if in.Avatar != "" {
		updateData["avatar"] = in.Avatar
	}
	if in.Email != "" {
		updateData["email"] = in.Email
	}

	// 更新用户信息
	err = l.svcCtx.DB.Model(&model.User{}).Where("id = ?", in.UserId).Updates(updateData).Error
	if err != nil {
		l.Logger.Errorf("更新用户信息失败: %v", err)
		return nil, status.Error(codes.Internal, "更新失败")
	}

	// 清除缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cacheKey := fmt.Sprintf("user:info:%d", in.UserId)
		l.svcCtx.Redis.Del(ctx, cacheKey)
	}()

	return &user.UpdateUserInfoResponse{
		Message: "用户信息更新成功",
	}, nil
}
