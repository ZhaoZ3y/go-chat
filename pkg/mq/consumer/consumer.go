package consumer

import (
	"IM/pkg/message"
	"IM/pkg/mq"
	"log"
)

type MessageConsumer struct {
	kafkaClient *mq.KafkaClient
	pushService PushService // WebSocket推送服务接口
}

type PushService interface {
	PushToUser(userID int64, message *message.PushMessage) error
	PushToGroup(groupID int64, message *message.PushMessage, excludeUserID int64) error
}

type PushMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewMessageConsumer(kafkaClient *mq.KafkaClient, pushService PushService) *MessageConsumer {
	return &MessageConsumer{
		kafkaClient: kafkaClient,
		pushService: pushService,
	}
}

func (c *MessageConsumer) Start() error {
	return c.kafkaClient.CreateConsumer(mq.TopicMessage, "im_message_consumer", c.handleMessage)
}

func (c *MessageConsumer) handleMessage(event *mq.MessageEvent) error {
	switch event.Type {
	case mq.EventNewMessage:
		return c.handleNewMessage(event)
	case mq.EventMessageRecall:
		return c.handleMessageRecall(event)
	case mq.EventMessageRead:
		return c.handleMessageRead(event)
	case mq.EventMessageDelete:
		return c.handleMessageDelete(event)
	default:
		log.Printf("未知消息类型: %s", event.Type)
		return nil
	}
}

func (c *MessageConsumer) handleNewMessage(event *mq.MessageEvent) error {
	pushMsg := &message.PushMessage{
		Type: "new_message",
		Data: map[string]interface{}{
			"message_id":   event.MessageID,
			"from_user_id": event.FromUserID,
			"to_user_id":   event.ToUserID,
			"group_id":     event.GroupID,
			"chat_type":    event.ChatType,
			"content":      event.Content,
			"message_type": event.MessageType,
			"extra":        event.Extra,
			"create_at":    event.CreateAt,
		},
	}

	if event.ChatType == 0 { // 私聊
		return c.pushService.PushToUser(event.ToUserID, pushMsg)
	} else { // 群聊
		return c.pushService.PushToGroup(event.GroupID, pushMsg, event.FromUserID)
	}
}

func (c *MessageConsumer) handleMessageRecall(event *mq.MessageEvent) error {
	pushMsg := &message.PushMessage{
		Type: "message_recall",
		Data: map[string]interface{}{
			"message_id":   event.MessageID,
			"from_user_id": event.FromUserID,
		},
	}

	if event.ChatType == 0 { // 私聊
		return c.pushService.PushToUser(event.ToUserID, pushMsg)
	} else { // 群聊
		return c.pushService.PushToGroup(event.GroupID, pushMsg, event.FromUserID)
	}
}

func (c *MessageConsumer) handleMessageRead(event *mq.MessageEvent) error {
	pushMsg := &message.PushMessage{
		Type: "message_read",
		Data: event.Data,
	}

	// 已读回执只推送给发送方
	return c.pushService.PushToUser(event.ToUserID, pushMsg)
}

func (c *MessageConsumer) handleMessageDelete(event *mq.MessageEvent) error {
	pushMsg := &message.PushMessage{
		Type: "message_delete",
		Data: map[string]interface{}{
			"message_id": event.MessageID,
		},
	}

	return c.pushService.PushToUser(event.FromUserID, pushMsg)
}
