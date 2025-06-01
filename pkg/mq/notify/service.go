package notify

import (
	"IM/pkg/mq"
	"fmt"
	"time"
)

// 通知类型常量
const (
	NotifyTypeJoinRequest   = "join_request"    // 加入群聊请求
	NotifyTypeLeaveGroup    = "leave_group"     // 退出群聊
	NotifyTypeInviteToGroup = "invite_to_group" // 邀请加入群聊
	NotifyTypeKickFromGroup = "kick_from_group" // 踢出群成员
	NotifyTypeDismissGroup  = "dismiss_group"   // 解散群聊
	NotifyTypeTransferGroup = "transfer_group"  // 转让群聊
	NotifyTypeMuteMember    = "mute_member"     // 禁言成员（群内通知）
)

// NotifyEvent 通知事件结构
type NotifyEvent struct {
	Type      string      `json:"type"`
	GroupID   int64       `json:"group_id"`
	GroupName string      `json:"group_name"`
	Data      interface{} `json:"data"`
}

// JoinRequestData 加入群聊请求通知数据
type JoinRequestData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Reason   string `json:"reason"`
}

// LeaveGroupData 退出群聊通知数据
type LeaveGroupData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// InviteToGroupData 邀请加入群聊通知数据
type InviteToGroupData struct {
	InviterID   int64    `json:"inviter_id"`
	InviterName string   `json:"inviter_name"`
	UserIDs     []int64  `json:"user_ids"`
	Usernames   []string `json:"usernames"`
}

// KickFromGroupData 踢出群成员通知数据
type KickFromGroupData struct {
	OperatorID   int64  `json:"operator_id"`
	OperatorName string `json:"operator_name"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
}

// DismissGroupData 解散群聊通知数据
type DismissGroupData struct {
	OwnerID   int64  `json:"owner_id"`
	OwnerName string `json:"owner_name"`
}

// TransferGroupData 转让群聊通知数据
type TransferGroupData struct {
	OldOwnerID   int64  `json:"old_owner_id"`
	OldOwnerName string `json:"old_owner_name"`
	NewOwnerID   int64  `json:"new_owner_id"`
	NewOwnerName string `json:"new_owner_name"`
}

// MuteMemberData 禁言成员通知数据
type MuteMemberData struct {
	OperatorID   int64  `json:"operator_id"`
	OperatorName string `json:"operator_name"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Duration     int64  `json:"duration"` // 0表示取消禁言
}

// NotifyService 通知服务接口
type NotifyService interface {
	// SendNotifyToAdmins 发送站外通知给群主和管理员
	SendNotifyToAdmins(event *NotifyEvent) error
	// SendNotifyToAllMembers 发送站外通知给所有群成员
	SendNotifyToAllMembers(event *NotifyEvent) error
	// SendGroupMessage 发送群内通知消息
	SendGroupMessage(event *NotifyEvent) error
}

// NotifyServiceImpl 通知服务实现
type NotifyServiceImpl struct {
	kafkaClient *mq.KafkaClient
}

func NewNotifyService(kafkaClient *mq.KafkaClient) NotifyService {
	return &NotifyServiceImpl{
		kafkaClient: kafkaClient,
	}
}

// SendNotifyToAdmins 发送站外通知给群主和管理员
func (s *NotifyServiceImpl) SendNotifyToAdmins(event *NotifyEvent) error {
	notifyEvent := &mq.MessageEvent{
		Type:    mq.EventNotifyAdmins,
		GroupID: event.GroupID,
		Data:    event,
	}
	return s.kafkaClient.SendMessage(mq.TopicNotify, notifyEvent)
}

// SendNotifyToAllMembers 发送站外通知给所有群成员
func (s *NotifyServiceImpl) SendNotifyToAllMembers(event *NotifyEvent) error {
	notifyEvent := &mq.MessageEvent{
		Type:    mq.EventNotifyAllMembers,
		GroupID: event.GroupID,
		Data:    event,
	}
	return s.kafkaClient.SendMessage(mq.TopicNotify, notifyEvent)
}

// SendGroupMessage 发送群内通知消息
func (s *NotifyServiceImpl) SendGroupMessage(event *NotifyEvent) error {
	// 构造群内消息
	content := s.buildNotifyMessage(event)

	messageEvent := &mq.MessageEvent{
		Type:        mq.EventNewMessage,
		GroupID:     event.GroupID,
		ChatType:    1, // 群聊
		Content:     content,
		MessageType: 1, // 文本消息
		CreateAt:    time.Now().Unix(),
	}

	return s.kafkaClient.SendMessage(mq.TopicMessage, messageEvent)
}

// buildNotifyMessage 构建通知消息内容
func (s *NotifyServiceImpl) buildNotifyMessage(event *NotifyEvent) string {
	switch event.Type {
	case NotifyTypeMuteMember:
		data := event.Data.(*MuteMemberData)
		if data.Duration == 0 {
			return fmt.Sprintf("%s 取消了 %s 的禁言", data.OperatorName, data.Username)
		}
		return fmt.Sprintf("%s 禁言了 %s", data.OperatorName, data.Username)

	case NotifyTypeJoinRequest:
		data := event.Data.(*JoinRequestData)
		return fmt.Sprintf("%s 加入了群聊", data.Username)

	case NotifyTypeInviteToGroup:
		data := event.Data.(*InviteToGroupData)
		if len(data.Usernames) == 1 {
			return fmt.Sprintf("%s 邀请 %s 加入了群聊", data.InviterName, data.Usernames[0])
		}
		return fmt.Sprintf("%s 邀请了 %d 人加入了群聊", data.InviterName, len(data.Usernames))

	case NotifyTypeKickFromGroup:
		data := event.Data.(*KickFromGroupData)
		return fmt.Sprintf("%s 将 %s 移出了群聊", data.OperatorName, data.Username)

	case NotifyTypeLeaveGroup:
		data := event.Data.(*LeaveGroupData)
		return fmt.Sprintf("%s 退出了群聊", data.Username)

	case NotifyTypeDismissGroup:
		data := event.Data.(*DismissGroupData)
		return fmt.Sprintf("群主 %s 解散了群聊", data.OwnerName)

	case NotifyTypeTransferGroup:
		data := event.Data.(*TransferGroupData)
		return fmt.Sprintf("%s 将群聊转让给了 %s", data.OldOwnerName, data.NewOwnerName)

	default:
		return "群组通知"
	}
}
