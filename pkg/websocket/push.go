package websocket

import (
	"IM/pkg/message"
	"encoding/json"
	"log"
)

// PushService WebSocket推送服务实现
type PushService struct {
	hub *Hub
}

func NewPushService(hub *Hub) *PushService {
	return &PushService{
		hub: hub,
	}
}

// PushToUser 推送消息给指定用户
func (p *PushService) PushToUser(userID int64, message *message.PushMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化推送消息失败: %v", err)
		return err
	}

	p.hub.SendToUser(userID, data)
	log.Printf("推送消息给用户 %d: %s", userID, message.Type)
	return nil
}

// PushToGroup 推送消息给群组成员（排除指定用户）
func (p *PushService) PushToGroup(groupID int64, message *message.PushMessage, excludeUserID int64) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化推送消息失败: %v", err)
		return err
	}

	p.hub.SendToGroup(groupID, data, excludeUserID)
	log.Printf("推送消息给群组 %d，排除用户 %d: %s", groupID, excludeUserID, message.Type)
	return nil
}

// PushToAll 推送消息给所有在线用户
func (p *PushService) PushToAll(message *message.PushMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化推送消息失败: %v", err)
		return err
	}

	p.hub.Broadcast(data)
	log.Printf("广播消息给所有用户: %s", message.Type)
	return nil
}
