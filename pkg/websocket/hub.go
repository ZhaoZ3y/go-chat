package websocket

import (
	"IM/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"sync"
)

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
	db         *gorm.DB
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

			logx.Infof("用户 %d 已连接WebSocket", client.UserID)

			// 用户上线后，推送离线期间的通知
			go h.pushOfflineNotifications(client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

			logx.Infof("用户 %d 已断开WebSocket连接", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.clients[userID]
	return exists
}

func (h *Hub) SendToUser(userID int64, message []byte) error {
	h.mu.RLock()
	client, exists := h.clients[userID]
	h.mu.RUnlock()

	if !exists {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	select {
	case client.Send <- message:
		return nil
	default:
		return fmt.Errorf("发送消息到用户 %d 失败", userID)
	}
}

// 推送离线期间的通知
func (h *Hub) pushOfflineNotifications(userID int64) {
	var notifications []model.Notifications

	// 查询用户离线期间的未读通知
	if err := h.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("create_at DESC").
		Limit(50).
		Find(&notifications).Error; err != nil {
		logx.Errorf("查询离线通知失败: %v", err)
		return
	}

	for _, notification := range notifications {
		notificationData, _ := json.Marshal(map[string]interface{}{
			"type":      "notification",
			"id":        notification.Id,
			"title":     notification.Title,
			"content":   notification.Content,
			"data":      notification.Data,
			"timestamp": notification.CreateAt,
		})

		if err := h.SendToUser(userID, notificationData); err != nil {
			logx.Errorf("推送离线通知失败: %v", err)
		}
	}

	// 标记为已读
	h.db.Model(&model.Notifications{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)
}
