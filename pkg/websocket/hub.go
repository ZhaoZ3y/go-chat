package websocket

import (
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
