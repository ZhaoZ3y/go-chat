package logic

import (
	"IM/pkg/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadFriendRequestCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnreadFriendRequestCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadFriendRequestCountLogic {
	return &GetUnreadFriendRequestCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取未读好友申请数量
func (l *GetUnreadFriendRequestCountLogic) GetUnreadFriendRequestCount(in *friend.GetUnreadFriendRequestCountRequest) (*friend.GetUnreadFriendRequestCountResponse, error) {
	var count int64
	err := l.svcCtx.DB.Model(&model.FriendRequests{}).
		Where("to_user_id = ? AND status = 1 AND is_read = false", in.UserId).
		Count(&count).Error

	if err != nil {
		return nil, status.Error(codes.Internal, "数据库查询失败")
	}

	return &friend.GetUnreadFriendRequestCountResponse{
		Count: int32(count),
	}, nil
}
