package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

// 在线用户连接管理
type WebSocketServer struct {
	upgrader websocket.Upgrader
	clients  map[int64]*websocket.Conn
	lock     sync.RWMutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[int64]*websocket.Conn),
	}
}

func (s *WebSocketServer) HandleWs(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	// 这里简化，实际可用 JWT 等鉴权
	// userID 转换
	var userID int64
	_, err := fmt.Sscan(userIDStr, &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logx.Errorf("upgrade websocket failed: %v", err)
		return
	}

	s.lock.Lock()
	s.clients[userID] = conn
	s.lock.Unlock()

	logx.Infof("user %d connected", userID)

	// 简单消息读取循环，防止连接关闭
	go func() {
		defer func() {
			s.lock.Lock()
			delete(s.clients, userID)
			s.lock.Unlock()
			conn.Close()
			logx.Infof("user %d disconnected", userID)
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

func (s *WebSocketServer) PushMessage(userID int64, msg string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	conn, ok := s.clients[userID]
	if !ok {
		return nil // 用户不在线，不报错
	}

	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
