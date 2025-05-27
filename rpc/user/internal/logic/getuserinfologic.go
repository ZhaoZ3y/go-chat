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

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	if in.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID无效")
	}

	// 先从Redis缓存获取
	cacheKey := fmt.Sprintf("user:info:%d", in.UserId)
	var userModel model.User

	// 从数据库查询
	err := l.svcCtx.DB.Where("id = ? AND deleted_at = 0", in.UserId).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 缓存用户信息到Redis (5分钟)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		l.svcCtx.Redis.Set(ctx, cacheKey, userModel, 5*time.Minute)
	}()

	return &user.GetUserInfoResponse{
		UserInfo: &user.User{
			Id:       userModel.Id,
			Username: userModel.Username,
			Email:    userModel.Email,
			Avatar:   userModel.Avatar,
			Phone:    userModel.Phone,
			Status:   int32(userModel.Status),
			CreateAt: userModel.CreateAt,
			UpdateAt: userModel.UpdateAt,
		},
	}, nil
}
