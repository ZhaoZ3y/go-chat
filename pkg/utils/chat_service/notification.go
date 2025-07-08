package chat_service

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"os"
	"time"
)

// NotificationService 封装了发送通知的逻辑。
type NotificationService struct {
	db    *gorm.DB
	kafka *mq.KafkaClient
}

// NewNotificationService 创建一个新的通知服务实例。
func NewNotificationService(db *gorm.DB, kafka *mq.KafkaClient) *NotificationService {
	return &NotificationService{
		db:    db,
		kafka: kafka,
	}
}

// SendSystemNotification 向单个用户发送系统通知。
func (s *NotificationService) SendSystemNotification(userID int64, title, content, extra string) error {
	serverNode := os.Getenv("SERVER_NODE")
	if serverNode == "" {
		serverNode = "notification-service"
	}

	event := &mq.NotificationEvent{
		Type:     mq.EventNotification,
		UserID:   userID,
		Title:    title,
		Content:  content,
		Extra:    extra,
		CreateAt: time.Now().Unix(),
	}

	// 将事件发送到 Kafka 的通知主题
	return s.kafka.SendMessage(mq.TopicNotification, event)
}

// SendBatchNotification 向多个用户批量发送相同的系统通知。
func (s *NotificationService) SendBatchNotification(userIDs []int64, title, content, extra string) error {
	// 循环调用单用户发送方法
	// 在高并发场景下，可以考虑优化为一次性生成所有事件再批量发送到 Kafka
	for _, userID := range userIDs {
		if err := s.SendSystemNotification(userID, title, content, extra); err != nil {
			// 记录失败的日志，但继续尝试发送给其他用户
			logx.Errorf("发送通知给用户 %d 失败: %v", userID, err)
		}
	}
	return nil
}

// SendGroupNotification 向一个群组的所有成员发送通知。
func (s *NotificationService) SendGroupNotification(groupID int64, title, content, extra string) error {
	// 1. 【重构】从数据库中查找群组的所有有效成员
	var members []model.GroupMembers
	// 使用 Where("deleted_at IS NULL") 来筛选当前在群里的成员
	if err := s.db.Model(&model.GroupMembers{}).Where("group_id = ? AND deleted_at IS NULL", groupID).Find(&members).Error; err != nil {
		logx.Errorf("查找群组 %d 成员失败: %v", groupID, err)
		return err
	}

	if len(members) == 0 {
		logx.Infof("群组 %d 没有有效成员，无需发送通知。", groupID)
		return nil
	}

	// 2. 收集所有成员的 UserID
	userIDs := make([]int64, 0, len(members))
	for _, member := range members {
		userIDs = append(userIDs, member.UserId)
	}

	// 3. 调用批量发送接口
	logx.Infof("正在向群组 %d 的 %d 位成员发送通知...", groupID, len(userIDs))
	return s.SendBatchNotification(userIDs, title, content, extra)
}
