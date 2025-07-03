package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadCountLogic {
	return &GetUnreadCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取总的未读数
func (l *GetUnreadCountLogic) GetUnreadCount(in *group.GetUnreadCountRequest) (*group.GetUnreadCountResponse, error) {
	var unreadCount int64
	err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&model.GroupNotification{}).
		Where("target_user_id = ? AND is_read = false", in.UserId).
		Count(&unreadCount).Error

	if err != nil {
		l.Logger.Errorf("failed to get unread count for user %d: %v", in.UserId, err)
		return nil, err
	}

	return &group.GetUnreadCountResponse{
		TotalUnreadCount: int32(unreadCount),
	}, nil
}
