package consumer

import (
	"IM/rpc/notify/internal/svc"
	"IM/rpc/notify/notification"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationConsumer struct {
	svcCtx *svc.ServiceContext
}

func NewNotificationConsumer(svcCtx *svc.ServiceContext) *NotificationConsumer {
	return &NotificationConsumer{svcCtx: svcCtx}
}

func (nc *NotificationConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var notification notification.NotificationMessage
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			logx.Errorf("unmarshal kafka message failed: %v", err)
			continue
		}

		logx.Infof("Kafka consumed notification: %+v", notification)
		// 这里可以调用RPC通知消费者接口或调用WebSocket推送

		sess.MarkMessage(msg, "")
	}
	return nil
}

func (nc *NotificationConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (nc *NotificationConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func RunConsumer(ctx context.Context, svcCtx *svc.ServiceContext) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(svcCtx.Config.Kafka.Brokers, svcCtx.Config.Kafka.ConsumerGroup, config)
	if err != nil {
		return err
	}
	defer client.Close()

	consumer := NewNotificationConsumer(svcCtx)

	for {
		if err := client.Consume(ctx, []string{svcCtx.Config.Kafka.ProducerTopic}, consumer); err != nil {
			logx.Errorf("Error from consumer: %v", err)
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
