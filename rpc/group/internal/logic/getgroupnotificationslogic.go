package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupNotificationsLogic {
	return &GetGroupNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组通知列表 (调用后，返回的通知在后端被标记为已读)
func (l *GetGroupNotificationsLogic) GetGroupNotifications(in *group.GetGroupNotificationsRequest) (*group.GetGroupNotificationsResponse, error) {
	var notifications []model.GroupNotification
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("target_user_id = ?", in.UserId).
		Order("timestamp DESC").
		Find(&notifications).Error

	if err != nil {
		l.Logger.Errorf("failed to get group notifications for user %d: %v", in.UserId, err)
		return nil, err
	}

	var pbNotifications []*group.GroupNotification
	var unreadIds []int64

	for _, notification := range notifications {
		pbNotification := &group.GroupNotification{
			Id:           notification.Id,
			Type:         group.NotificationType(notification.Type),
			GroupId:      notification.GroupId,
			OperatorId:   notification.OperatorId,
			TargetUserId: notification.TargetUserId,
			Message:      notification.Message,
			Timestamp:    notification.CreateAt,
			IsRead:       notification.IsRead,
		}
		pbNotifications = append(pbNotifications, pbNotification)

		// 收集未读通知的ID
		if !notification.IsRead {
			unreadIds = append(unreadIds, notification.Id)
		}
	}

	// 将未读通知标记为已读
	if len(unreadIds) > 0 {
		err = l.svcCtx.DB.WithContext(l.ctx).
			Model(&model.GroupNotification{}).
			Where("id IN ?", unreadIds).
			Update("is_read", true).Error

		if err != nil {
			l.Logger.Errorf("failed to mark notifications as read: %v", err)
		} else {
			for _, pbNotification := range pbNotifications {
				pbNotification.IsRead = true
			}
		}
	}

	return &group.GetGroupNotificationsResponse{
		Notifications: pbNotifications,
	}, nil
}
