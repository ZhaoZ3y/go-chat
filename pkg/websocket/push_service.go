package websocket

import (
	"IM/pkg/model"
	"IM/rpc/notify/notification"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type PushService struct {
	hub *Hub
	db  *gorm.DB
}

func NewPushService(hub *Hub) *PushService {
	return &PushService{
		hub: hub,
		db:  hub.db,
	}
}

// 处理通知消息
func (ps *PushService) HandleNotification(ctx context.Context, notifyMsg *notification.NotificationMessage) error {
	// 先存储通知到数据库
	dbNotification := &model.Notifications{
		UserId:   notifyMsg.UserId,
		Type:     int8(notifyMsg.Type),
		Title:    ps.getNotificationTitle(notifyMsg.Type),
		Content:  notifyMsg.Content,
		Data:     notifyMsg.Extra,
		IsRead:   false,
		CreateAt: time.Now().Unix(),
	}

	if err := ps.db.Create(dbNotification).Error; err != nil {
		logx.Errorf("保存通知到数据库失败: %v", err)
		return err
	}

	// 如果用户在线，立即推送
	if ps.hub.IsUserOnline(notifyMsg.UserId) {
		return ps.pushToOnlineUser(notifyMsg, dbNotification.Id)
	}

	logx.Infof("用户 %d 不在线，通知已保存到数据库", notifyMsg.UserId)
	return nil
}

// 推送给在线用户
func (ps *PushService) pushToOnlineUser(notifyMsg *notification.NotificationMessage, notificationId int64) error {
	pushData := map[string]interface{}{
		"type":              "notification",
		"id":                notificationId,
		"title":             ps.getNotificationTitle(notifyMsg.Type),
		"content":           notifyMsg.Content,
		"notification_type": int(notifyMsg.Type),
		"group_id":          notifyMsg.GroupId,
		"operator_id":       notifyMsg.OperatorUserId,
		"data":              notifyMsg.Extra,
		"timestamp":         notifyMsg.Timestamp,
	}

	messageData, err := json.Marshal(pushData)
	if err != nil {
		return err
	}

	return ps.hub.SendToUser(notifyMsg.UserId, messageData)
}

// 获取通知标题
func (ps *PushService) getNotificationTitle(notifyType notification.NotificationType) string {
	switch notifyType {
	case notification.NotificationType_FRIEND_REQUEST:
		return "好友申请"
	case notification.NotificationType_GROUP_JOIN_REQUEST:
		return "群聊申请"
	case notification.NotificationType_GROUP_JOINED_BROADCAST:
		return "群聊通知"
	case notification.NotificationType_GROUP_INVITE_JOIN_REQUEST:
		return "群聊邀请"
	case notification.NotificationType_GROUP_HANDLE_JOIN_REQUEST:
		return "申请处理结果"
	case notification.NotificationType_GROUP_BEINVITED_JOINED_REQUEST:
		return "被邀请入群"
	case notification.NotificationType_GROUP_MEMBER_KICKED:
		return "群成员变动"
	case notification.NotificationType_GROUP_MEMBER_QUIT:
		return "群成员退出"
	case notification.NotificationType_GROUP_TRANSFERRED:
		return "群聊转让"
	case notification.NotificationType_GROUP_ADMIN_HANDLE:
		return "群管理通知"
	case notification.NotificationType_GROUP_DISMISSED:
		return "群聊解散"
	default:
		return "系统通知"
	}
}
