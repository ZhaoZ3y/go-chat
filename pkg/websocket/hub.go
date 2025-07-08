package websocket

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	_const "IM/pkg/utils/const"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Client struct {
	UserID     int64
	Conn       *websocket.Conn
	Send       chan []byte
	ServerNode string
}

type Hub struct {
	clients      map[int64]*Client
	register     chan *Client
	unregister   chan *Client
	broadcast    chan []byte
	mu           sync.RWMutex
	db           *gorm.DB
	kafka        *mq.KafkaClient
	redisCli     *redis.Client // 用于管理在线状态的 Redis 客户端
	serverNode   string
	pingInterval time.Duration
	pongWait     time.Duration
	ctx          context.Context
	cancel       context.CancelFunc
}

type WSMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// NewHub 现在需要接收一个 redis.Client
func NewHub(db *gorm.DB, kafka *mq.KafkaClient, redisCli *redis.Client, serverNode string) *Hub {
	ctx, cancel := context.WithCancel(context.Background())

	return &Hub{
		clients:      make(map[int64]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan []byte),
		db:           db,
		kafka:        kafka,
		redisCli:     redisCli,
		serverNode:   serverNode,
		pingInterval: 30 * time.Second,
		pongWait:     60 * time.Second, // 必须大于 pingInterval
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (h *Hub) Run() {
	// 启动 Kafka 消息消费者
	go h.startMessageConsumer()

	for {
		select {
		case <-h.ctx.Done():
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

			// 在 Redis 中更新用户在线状态
			h.updateUserStatusInRedis(client.UserID, 1)
			// 向 Kafka 发布用户上线事件 (供其他微服务使用)
			h.publishUserStatusEvent(client.UserID, mq.EventUserOnline)
			// 推送离线消息
			go h.pushOfflineMessages(client.UserID)

			logx.Infof("用户 %d 已连接至 WebSocket 节点 %s", client.UserID, h.serverNode)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)

				// 在 Redis 中更新用户离线状态
				h.updateUserStatusInRedis(client.UserID, 0)
				// 向 Kafka 发布用户下线事件
				h.publishUserStatusEvent(client.UserID, mq.EventUserOffline)
			}
			h.mu.Unlock()

			logx.Infof("用户 %d 已从 WebSocket 节点 %s 断开", client.UserID, h.serverNode)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// 如果发送缓冲区已满，则认为客户端出现问题，断开其连接
					h.unregister <- client
				}
			}
			h.mu.RUnlock()
		}
	}
}

// updateUserStatusInRedis 负责在 Redis 中处理所有在线状态更新。
// 它会设置一个 HASH，包含状态、节点和时间戳，并设置一个过期时间。
func (h *Hub) updateUserStatusInRedis(userID int64, status int8) {
	ctx := context.Background()
	userKey := fmt.Sprintf("im:status:%d", userID)
	nodeKey := fmt.Sprintf("im:node_users:%s", h.serverNode)
	now := time.Now().Unix()

	pipe := h.redisCli.Pipeline()

	if status == 1 { // 用户上线
		// HASH: im:status:{userID} -> {status: 1, server_node: "node-1", last_seen: 167...}
		pipe.HSet(ctx, userKey, "status", 1, "server_node", h.serverNode, "last_seen", now)
		// 每次有活动（连接、心跳）时都刷新 TTL
		pipe.Expire(ctx, userKey, h.pongWait+10*time.Second)
		// 将用户添加到此服务器节点的在线用户集合中
		pipe.SAdd(ctx, nodeKey, userID)
	} else { // 用户下线
		// 只更新状态为离线，让 TTL 机制来处理键的最终移除
		pipe.HSet(ctx, userKey, "status", 0, "last_seen", now)
		// 从此服务器节点的集合中移除用户
		pipe.SRem(ctx, nodeKey, userID)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		logx.Errorf("在 Redis 中更新用户 %d 状态为 %d 失败: %v", userID, status, err)
	}
}

// publishUserStatusEvent 向 Kafka 发送一个事件。其他微服务可能对用户的上下线事件感兴趣。
func (h *Hub) publishUserStatusEvent(userID int64, eventType string) {
	status := int8(1)
	if eventType == mq.EventUserOffline {
		status = 0
	}

	event := &mq.UserStatusEvent{
		Type:      eventType,
		UserID:    userID,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}

	if err := h.kafka.SendMessage(mq.TopicUserStatus, event); err != nil {
		logx.Errorf("向 Kafka 发布用户状态事件失败: %v", err)
	}
}

func (h *Hub) pushOfflineMessages(userID int64) {
	var offlineMessages []model.OfflineMessages
	// 从数据库中获取未推送的消息
	h.db.Where("user_id = ? AND status = 0", userID).Find(&offlineMessages)

	if len(offlineMessages) == 0 {
		return
	}

	logx.Infof("正在向用户 %d 推送 %d 条离线消息", len(offlineMessages), userID)
	for _, offlineMsg := range offlineMessages {
		// 离线消息的内容就是完整的 WebSocket 消息 payload，无需再次查询原始消息。
		if h.SendToUser(userID, []byte(offlineMsg.Content)) {
			// 只有成功发送到客户端的 channel 后，才标记为已推送
			h.db.Model(&offlineMsg).Update("status", 1)
		}
	}
}

func (h *Hub) startMessageConsumer() {
	// 消费普通消息
	h.kafka.ConsumeMessages(mq.TopicMessage, func(data []byte) error {
		var event mq.MessageEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		return h.handleMessageEvent(&event)
	})

	// 消费通知
	h.kafka.ConsumeMessages(mq.TopicNotification, func(data []byte) error {
		var event mq.NotificationEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		return h.handleNotificationEvent(&event)
	})
}

func (h *Hub) handleMessageEvent(event *mq.MessageEvent) error {
	switch event.Type {
	case mq.EventNewMessage:
		return h.handleNewMessage(event)
	case mq.EventMessageRead:
		return h.handleMessageRead(event)
	case mq.EventMessageRecall:
		return h.handleMessageRecall(event)
	}
	return nil
}

func (h *Hub) handleNewMessage(event *mq.MessageEvent) error {
	wsMsg := WSMessage{
		Type:      "new_message",
		Data:      event,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	if event.ChatType == 0 { // 私聊
		if !h.SendToUser(event.ToUserID, data) {
			h.saveOfflineMessage(event.ToUserID, event.MessageID, 0, string(data))
		}
	} else { // 群聊
		var members []model.GroupMembers
		h.db.Where("group_id = ? AND deleted_at IS NULL", event.GroupID).Find(&members)

		for _, member := range members {
			if member.UserId == event.FromUserID {
				continue // 不发给自己
			}
			if !h.SendToUser(member.UserId, data) {
				// 用户不在此节点在线，保存为离线消息
				h.saveOfflineMessage(member.UserId, event.MessageID, 0, string(data))
			}
		}
	}

	return nil
}

func (h *Hub) handleMessageRead(event *mq.MessageEvent) error {
	h.db.Model(&model.Messages{}).
		Where("id = ?", event.MessageID).
		Update("is_read", 1)

	// 通知原始发送者消息已被阅读
	wsMsg := WSMessage{
		Type:      "message_read",
		Data:      event,
		Timestamp: time.Now().Unix(),
	}

	if data, err := json.Marshal(wsMsg); err == nil {
		h.SendToUser(event.FromUserID, data)
	}

	return nil
}

func (h *Hub) handleNotificationEvent(event *mq.NotificationEvent) error {
	wsMsg := WSMessage{
		Type:      "notification",
		Data:      event,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	if !h.SendToUser(event.UserID, data) {
		// 注意：系统通知的 messageID 为 0
		h.saveOfflineMessage(event.UserID, 0, 1, string(data))
	}

	return nil
}

func (h *Hub) handleMessageRecall(event *mq.MessageEvent) error {
	wsMsg := WSMessage{
		Type:      "message_recalled", // 客户端需要识别这个类型
		Data:      event,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		logx.Errorf("序列化消息撤回通知失败: %v", err)
		return err
	}

	if event.ChatType == _const.ChatTypePrivate { // 私聊
		// 通知接收者
		h.SendToUser(event.ToUserID, data)
		// 同时也要通知发送者（可能在其他设备上登录），使其同步状态
		h.SendToUser(event.FromUserID, data)
	} else { // 群聊
		// 获取群组所有成员
		var members []model.GroupMembers
		h.db.Where("group_id = ? AND deleted_at IS NULL", event.GroupID).Find(&members)

		// 向群内所有成员广播撤回通知
		for _, member := range members {
			h.SendToUser(member.UserId, data)
		}
	}

	return nil
}

func (h *Hub) saveOfflineMessage(userID, messageID int64, msgType int8, content string) {
	offlineMsg := &model.OfflineMessages{
		UserId:    userID,
		MessageId: messageID,
		Type:      msgType,
		Content:   content,
		Status:    0, // 未推送
	}
	if err := h.db.Create(offlineMsg).Error; err != nil {
		logx.Errorf("为用户 %d 保存离线消息失败: %v", userID, err)
	}
}

// IsUserOnline 检查用户是否连接到
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.clients[userID]
	return exists
}

// SendToUser 如果用户连接到此 Hub 实例，则向其发送消息
func (h *Hub) SendToUser(userID int64, message []byte) bool {
	h.mu.RLock()
	client, exists := h.clients[userID]
	h.mu.RUnlock()

	if !exists {
		return false
	}

	select {
	case client.Send <- message:
		return true
	default:
		logx.Errorf("向用户 %d 发送消息失败: 发送 channel 已满", userID)
		go func() {
			h.unregister <- client
		}()
		return false
	}
}

func (h *Hub) GetOnlineUsers() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]int64, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

func (h *Hub) Close() {
	h.cancel()
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, client := range h.clients {
		close(client.Send)
		client.Conn.Close()
	}
}
