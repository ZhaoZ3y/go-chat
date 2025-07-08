package websocket

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 为简单起见，允许所有来源
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// PingMessage 是从服务器发送到客户端的
type PingMessage struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	userIDStr := c.GetString("userID")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 user_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logx.Errorf("为用户 %d 升级 websocket 失败: %v", userID, err)
		return
	}

	client := &Client{
		UserID:     userID,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		ServerNode: h.serverNode,
	}

	h.register <- client

	// 为此客户端启动读写 goroutine
	go h.writePump(client)
	go h.readPump(client)
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(512) // 设置最大消息大小
	client.Conn.SetReadDeadline(time.Now().Add(h.pongWait))

	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(h.pongWait))
		h.updateUserStatusInRedis(client.UserID, 1)
		return nil
	})

	for {
		_, messageBytes, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("用户 %d 的 Websocket 错误: %v", client.UserID, err)
			}
			break // 出现错误时退出循环
		}

		// 处理来自客户端的传入消息
		h.handleClientMessage(client, messageBytes)
	}
}

func (h *Hub) writePump(client *Client) {
	ticker := time.NewTicker(h.pingInterval)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub 关闭了 channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logx.Errorf("用户 %d 的 Websocket 写入错误: %v", client.UserID, err)
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			// 发送一个 JSON 格式的 ping 消息，而不是原始的 ping 帧，这样更灵活
			pingMsg := PingMessage{
				Type:      "ping",
				Timestamp: time.Now().Unix(),
			}
			data, _ := json.Marshal(pingMsg)
			if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				logx.Errorf("用户 %d 的 Ping 错误: %v", client.UserID, err)
				return
			}
		}
	}
}

func (h *Hub) handleClientMessage(client *Client, message []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		logx.Errorf("解析客户端消息失败: %v", err)
		return
	}

	switch wsMsg.Type {
	case "pong":
		client.Conn.SetReadDeadline(time.Now().Add(h.pongWait))
		h.updateUserStatusInRedis(client.UserID, 1)
	case "typing":
		h.handleTypingStatus(client, wsMsg.Data)
	default:
		logx.Infof("来自客户端 %d 的未知消息类型: %s", client.UserID, wsMsg.Type)
	}
}

func (h *Hub) handleMessageReadAck(userID, messageID int64) {
	var message model.Messages
	// 查找消息以获取发送者的ID (FromUserId)
	if err := h.db.First(&message, "id = ? AND to_user_id = ?", messageID, userID).Error; err != nil {
		// 消息未找到或不属于此用户
		return
	}

	// 创建一个消息已读事件并发送到 Kafka
	readEvent := &mq.MessageEvent{
		Type:       mq.EventMessageRead,
		MessageID:  messageID,
		FromUserID: message.FromUserId, // 原始发送者
		ToUserID:   userID,             // 阅读消息的人
		ChatType:   message.ChatType,
		GroupID:    message.GroupId,
		CreateAt:   time.Now().Unix(),
	}

	h.kafka.SendMessage(mq.TopicMessage, readEvent)
}

func (h *Hub) handleTypingStatus(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	dataMap["from_user_id"] = client.UserID

	wsMsg := WSMessage{
		Type:      "typing",
		Data:      dataMap,
		Timestamp: time.Now().Unix(),
	}
	msgData, err := json.Marshal(wsMsg)
	if err != nil {
		return
	}

	if chatType, ok := dataMap["chat_type"].(float64); ok {
		if int(chatType) == 0 { // 私聊
			if toUserID, ok := dataMap["to_user_id"].(float64); ok {
				h.SendToUser(int64(toUserID), msgData)
			}
		} else { // 群聊
			if groupID, ok := dataMap["group_id"].(float64); ok {
				h.broadcastToGroup(int64(groupID), client.UserID, msgData)
			}
		}
	}
}

func (h *Hub) broadcastToGroup(groupID, excludeUserID int64, message []byte) {
	var members []model.GroupMembers
	h.db.Where("group_id = ? AND deleted_at IS NULL", groupID).Find(&members)

	for _, member := range members {
		if member.UserId == excludeUserID {
			continue
		}
		h.SendToUser(member.UserId, message)
	}
}
