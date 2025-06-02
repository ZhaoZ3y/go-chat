package consumer

import (
	"IM/pkg/websocket"
	"IM/rpc/notify/notification"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type NotifyConsumer struct {
	brokers     []string
	pushService *websocket.PushService
	db          *gorm.DB
}

func NewNotifyConsumer(brokers []string, pushService *websocket.PushService, db *gorm.DB) *NotifyConsumer {
	return &NotifyConsumer{
		brokers:     brokers,
		pushService: pushService,
		db:          db,
	}
}

func (nc *NotifyConsumer) Start() error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(
		nc.brokers,
		"notification_consumer_group",
		config,
	)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	ctx := context.Background()
	for {
		if err := consumerGroup.Consume(ctx, []string{"notification_topic"}, nc); err != nil {
			logx.Errorf("通知消费者错误: %v", err)
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (nc *NotifyConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (nc *NotifyConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (nc *NotifyConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var notifyMsg notification.NotificationMessage
		if err := json.Unmarshal(message.Value, &notifyMsg); err != nil {
			logx.Errorf("解析通知消息失败: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		logx.Infof("收到通知消息: UserID=%d, Type=%v, Content=%s",
			notifyMsg.UserId, notifyMsg.Type, notifyMsg.Content)

		// 处理通知消息
		if err := nc.pushService.HandleNotification(context.Background(), &notifyMsg); err != nil {
			logx.Errorf("处理通知消息失败: %v", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}
