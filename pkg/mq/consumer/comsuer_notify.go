package consumer

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"IM/pkg/mq/notify"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type NotifyConsumer struct {
	kafkaClient *mq.KafkaClient
	pushService PushService
	db          *gorm.DB
}

func NewNotifyConsumer(kafkaClient *mq.KafkaClient, pushService PushService, db *gorm.DB) *NotifyConsumer {
	return &NotifyConsumer{
		kafkaClient: kafkaClient,
		pushService: pushService,
		db:          db,
	}
}

func (c *NotifyConsumer) Start() error {
	return c.kafkaClient.CreateConsumer(mq.TopicNotify, "im_notify_consumer", c.handleNotify)
}

func (c *NotifyConsumer) handleNotify(event *mq.MessageEvent) error {
	switch event.Type {
	case mq.EventNotifyAdmins:
		return c.handleNotifyAdmins(event)
	case mq.EventNotifyAllMembers:
		return c.handleNotifyAllMembers(event)
	default:
		log.Printf("未知通知类型: %s", event.Type)
		return nil
	}
}

// 处理通知群主和管理员
func (c *NotifyConsumer) handleNotifyAdmins(event *mq.MessageEvent) error {
	// 解析通知事件
	notifyEvent := &notify.NotifyEvent{}
	if eventData, ok := event.Data.(map[string]interface{}); ok {
		data, _ := json.Marshal(eventData)
		json.Unmarshal(data, notifyEvent)
	}

	// 获取群主和管理员列表
	var admins []model.GroupMembers
	if err := c.db.Where("group_id = ? AND role IN (1,2)", event.GroupID).Find(&admins).Error; err != nil {
		log.Printf("获取群管理员失败: %v", err)
		return err
	}

	// 构造推送消息
	pushMsg := &PushMessage{
		Type: "group_notify",
		Data: map[string]interface{}{
			"notify_type": notifyEvent.Type,
			"group_id":    notifyEvent.GroupID,
			"group_name":  notifyEvent.GroupName,
			"message":     c.buildNotifyMessage(notifyEvent),
			"data":        notifyEvent.Data,
			"timestamp":   time.Now().Unix(),
		},
	}

	// 推送给每个管理员
	for _, admin := range admins {
		if err := c.pushService.PushToUser(admin.UserId, pushMsg); err != nil {
			log.Printf("推送通知给管理员 %d 失败: %v", admin.UserId, err)
		}
	}

	return nil
}

// 处理通知所有群成员
func (c *NotifyConsumer) handleNotifyAllMembers(event *mq.MessageEvent) error {
	// 解析通知事件
	notifyEvent := &notify.NotifyEvent{}
	if eventData, ok := event.Data.(map[string]interface{}); ok {
		data, _ := json.Marshal(eventData)
		json.Unmarshal(data, notifyEvent)
	}

	// 获取所有群成员列表
	var members []model.GroupMembers
	if err := c.db.Where("group_id = ?", event.GroupID).Find(&members).Error; err != nil {
		log.Printf("获取群成员失败: %v", err)
		return err
	}

	// 构造推送消息
	pushMsg := &PushMessage{
		Type: "group_notify",
		Data: map[string]interface{}{
			"notify_type": notifyEvent.Type,
			"group_id":    notifyEvent.GroupID,
			"group_name":  notifyEvent.GroupName,
			"message":     c.buildNotifyMessage(notifyEvent),
			"data":        notifyEvent.Data,
			"timestamp":   time.Now().Unix(),
		},
	}

	// 推送给每个群成员
	for _, member := range members {
		if err := c.pushService.PushToUser(member.UserId, pushMsg); err != nil {
			log.Printf("推送通知给成员 %d 失败: %v", member.UserId, err)
		}
	}

	return nil
}

// 构建通知消息
func (c *NotifyConsumer) buildNotifyMessage(event *notify.NotifyEvent) string {
	switch event.Type {
	case notify.NotifyTypeJoinRequest:
		data := event.Data.(map[string]interface{})
		username := data["username"].(string)
		return fmt.Sprintf("%s 申请加入群聊", username)

	case notify.NotifyTypeLeaveGroup:
		data := event.Data.(map[string]interface{})
		username := data["username"].(string)
		return fmt.Sprintf("%s 退出了群聊", username)

	case notify.NotifyTypeInviteToGroup:
		data := event.Data.(map[string]interface{})
		inviterName := data["inviter_name"].(string)
		usernames := data["usernames"].([]interface{})
		names := make([]string, len(usernames))
		for i, name := range usernames {
			names[i] = name.(string)
		}
		return fmt.Sprintf("%s 邀请了 %v 加入群聊", inviterName, names)

	case notify.NotifyTypeKickFromGroup:
		data := event.Data.(map[string]interface{})
		operatorName := data["operator_name"].(string)
		username := data["username"].(string)
		return fmt.Sprintf("%s 将 %s 踢出了群聊", operatorName, username)

	case notify.NotifyTypeDismissGroup:
		data := event.Data.(map[string]interface{})
		ownerName := data["owner_name"].(string)
		return fmt.Sprintf("群主 %s 解散了群聊", ownerName)

	case notify.NotifyTypeTransferGroup:
		data := event.Data.(map[string]interface{})
		oldOwnerName := data["old_owner_name"].(string)
		newOwnerName := data["new_owner_name"].(string)
		return fmt.Sprintf("%s 将群聊转让给了 %s", oldOwnerName, newOwnerName)

	default:
		return "群组通知"
	}
}
