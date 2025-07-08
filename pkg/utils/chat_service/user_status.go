package chat_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type UserStatus struct {
	UserID     int64  `json:"user_id"`
	Status     int8   `json:"status"`      // 0: 离线, 1: 在线
	LastSeen   int64  `json:"last_seen"`   // Unix 时间戳
	ServerNode string `json:"server_node"` // 用户连接到的服务器节点
}

// UserStatusService 现在只依赖于 Redis。
type UserStatusService struct {
	redisCli *redis.Client
}

// NewUserStatusService 创建一个带有 Redis 客户端的新服务。
func NewUserStatusService(redisCli *redis.Client) *UserStatusService {
	return &UserStatusService{
		redisCli: redisCli,
	}
}

// GetUserStatus 从 Redis 中获取单个用户的在线状态。
func (s *UserStatusService) GetUserStatus(userID int64) (*UserStatus, error) {
	key := fmt.Sprintf("im:status:%d", userID)
	data, err := s.redisCli.HGetAll(context.Background(), key).Result()

	// 如果键不存在 (redis.Nil) 或为空，则认为用户离线。
	if errors.Is(err, redis.Nil) || len(data) == 0 {
		return &UserStatus{UserID: userID, Status: 0}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("从 redis 获取用户状态失败: %w", err)
	}

	status, _ := strconv.ParseInt(data["status"], 10, 8)
	lastSeen, _ := strconv.ParseInt(data["last_seen"], 10, 64)

	return &UserStatus{
		UserID:     userID,
		Status:     int8(status),
		LastSeen:   lastSeen,
		ServerNode: data["server_node"],
	}, nil
}

// GetBatchUserStatus 使用 Redis pipeline 高效地获取多个用户的状态。
func (s *UserStatusService) GetBatchUserStatus(userIDs []int64) (map[int64]*UserStatus, error) {
	if len(userIDs) == 0 {
		return make(map[int64]*UserStatus), nil
	}

	ctx := context.Background()
	pipe := s.redisCli.Pipeline()
	cmds := make(map[int64]*redis.MapStringStringCmd)

	for _, userID := range userIDs {
		key := fmt.Sprintf("im:status:%d", userID)
		cmds[userID] = pipe.HGetAll(ctx, key)
	}

	// 执行 pipeline 会一次性将所有命令发送到 Redis。
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("redis pipeline 执行失败: %w", err)
	}

	result := make(map[int64]*UserStatus, len(userIDs))
	for userID, cmd := range cmds {
		data, err := cmd.Result()
		// 处理特定用户的键不存在的情况。
		if err == redis.Nil || len(data) == 0 {
			result[userID] = &UserStatus{UserID: userID, Status: 0}
			continue
		}

		status, _ := strconv.ParseInt(data["status"], 10, 8)
		lastSeen, _ := strconv.ParseInt(data["last_seen"], 10, 64)

		result[userID] = &UserStatus{
			UserID:     userID,
			Status:     int8(status),
			LastSeen:   lastSeen,
			ServerNode: data["server_node"],
		}
	}

	return result, nil
}
