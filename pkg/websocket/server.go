package websocket

import (
	"IM/pkg/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients      map[int64]*Client
	groupMembers map[int64][]int64 // groupID -> []userID
	register     chan *Client
	unregister   chan *Client
	broadcast    chan []byte
	mutex        sync.RWMutex
	db           *gorm.DB
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID int64
	mutex  sync.Mutex
}

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		clients:      make(map[int64]*Client),
		groupMembers: make(map[int64][]int64),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan []byte),
		db:           db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.loadUserGroups(client.userID)
			h.mutex.Unlock()
			log.Printf("用户 %d 已连接，当前在线用户数: %d", client.userID, len(h.clients))

			// 发送连接成功消息
			welcomeMsg := &Message{
				Type: "connection_success",
				Data: map[string]interface{}{
					"user_id": client.userID,
					"message": "WebSocket连接成功",
				},
				Timestamp: time.Now().Unix(),
			}
			if data, err := json.Marshal(welcomeMsg); err == nil {
				client.send <- data
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				h.removeUserFromGroups(client.userID)
			}
			h.mutex.Unlock()
			log.Printf("用户 %d 已断开连接，当前在线用户数: %d", client.userID, len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// loadUserGroups 加载用户所属的群组
func (h *Hub) loadUserGroups(userID int64) {
	var groupMembers []model.GroupMembers
	if err := h.db.Where("user_id = ?", userID).Find(&groupMembers).Error; err != nil {
		log.Printf("加载用户 %d 的群组失败: %v", userID, err)
		return
	}

	for _, member := range groupMembers {
		if _, exists := h.groupMembers[member.GroupId]; !exists {
			h.groupMembers[member.GroupId] = make([]int64, 0)
		}

		// 检查用户是否已在群组列表中
		found := false
		for _, uid := range h.groupMembers[member.GroupId] {
			if uid == userID {
				found = true
				break
			}
		}

		if !found {
			h.groupMembers[member.GroupId] = append(h.groupMembers[member.GroupId], userID)
		}
	}
}

// removeUserFromGroups 从所有群组中移除用户
func (h *Hub) removeUserFromGroups(userID int64) {
	for groupID, members := range h.groupMembers {
		for i, uid := range members {
			if uid == userID {
				h.groupMembers[groupID] = append(members[:i], members[i+1:]...)
				break
			}
		}

		// 如果群组没有在线成员，清空群组
		if len(h.groupMembers[groupID]) == 0 {
			delete(h.groupMembers, groupID)
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, message []byte) {
	h.mutex.RLock()
	client, ok := h.clients[userID]
	h.mutex.RUnlock()

	if ok {
		select {
		case client.send <- message:
			log.Printf("消息已发送给用户 %d", userID)
		default:
			h.mutex.Lock()
			close(client.send)
			delete(h.clients, userID)
			h.removeUserFromGroups(userID)
			h.mutex.Unlock()
			log.Printf("用户 %d 连接已关闭", userID)
		}
	} else {
		log.Printf("用户 %d 不在线，消息发送失败", userID)
	}
}

// SendToGroup 发送消息给群组成员（排除指定用户）
func (h *Hub) SendToGroup(groupID int64, message []byte, excludeUserID int64) {
	h.mutex.RLock()
	members, ok := h.groupMembers[groupID]
	if !ok {
		h.mutex.RUnlock()
		log.Printf("群组 %d 没有在线成员", groupID)
		return
	}

	sentCount := 0
	for _, userID := range members {
		if userID == excludeUserID {
			continue
		}

		if client, exists := h.clients[userID]; exists {
			select {
			case client.send <- message:
				sentCount++
			default:
				close(client.send)
				delete(h.clients, userID)
			}
		}
	}
	h.mutex.RUnlock()

	log.Printf("群组消息已发送给 %d 个成员，群组ID: %d", sentCount, groupID)
}

// Broadcast 广播消息给所有在线用户
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// GetOnlineUserCount 获取在线用户数量
func (h *Hub) GetOnlineUserCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetOnlineUsers 获取在线用户列表
func (h *Hub) GetOnlineUsers() []int64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]int64, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// 设置读取超时
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 处理客户端发送的消息（心跳、状态等）
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil {
			c.handleClientMessage(&msg)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("发送消息失败: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleClientMessage 处理客户端消息
func (c *Client) handleClientMessage(msg *Message) {
	switch msg.Type {
	case "ping":
		// 响应心跳
		pongMsg := &Message{
			Type:      "pong",
			Timestamp: time.Now().Unix(),
		}
		if data, err := json.Marshal(pongMsg); err == nil {
			select {
			case c.send <- data:
			default:
			}
		}
	case "join_group":
		// 处理加入群组
		if data, ok := msg.Data.(map[string]interface{}); ok {
			if groupIDFloat, exists := data["group_id"]; exists {
				groupID := int64(groupIDFloat.(float64))
				c.hub.mutex.Lock()
				if _, exists := c.hub.groupMembers[groupID]; !exists {
					c.hub.groupMembers[groupID] = make([]int64, 0)
				}

				// 检查用户是否已在群组中
				found := false
				for _, uid := range c.hub.groupMembers[groupID] {
					if uid == c.userID {
						found = true
						break
					}
				}

				if !found {
					c.hub.groupMembers[groupID] = append(c.hub.groupMembers[groupID], c.userID)
				}
				c.hub.mutex.Unlock()

				log.Printf("用户 %d 加入群组 %d", c.userID, groupID)
			}
		}
	case "leave_group":
		// 处理退出群组
		if data, ok := msg.Data.(map[string]interface{}); ok {
			if groupIDFloat, exists := data["group_id"]; exists {
				groupID := int64(groupIDFloat.(float64))
				c.hub.mutex.Lock()
				if members, exists := c.hub.groupMembers[groupID]; exists {
					for i, uid := range members {
						if uid == c.userID {
							c.hub.groupMembers[groupID] = append(members[:i], members[i+1:]...)
							break
						}
					}
				}
				c.hub.mutex.Unlock()

				log.Printf("用户 %d 退出群组 %d", c.userID, groupID)
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 从请求中获取用户ID（应该从JWT token中解析）
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		log.Println("缺少用户ID参数")
		conn.Close()
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Printf("无效的用户ID: %s", userIDStr)
		conn.Close()
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}

	client.hub.register <- client

	// 启动goroutines
	go client.writePump()
	go client.readPump()
}
