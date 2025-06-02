package consumer

import (
	"IM/pkg/websocket"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageConsumer struct {
	brokers     []string
	pushService *websocket.PushService
}

func NewMessageConsumer(brokers []string, pushService *websocket.PushService) *MessageConsumer {
	return &MessageConsumer{
		brokers:     brokers,
		pushService: pushService,
	}
}

func (mc *MessageConsumer) Start() error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(
		mc.brokers,
		"message_consumer_group",
		config,
	)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	ctx := context.Background()
	for {
		if err := consumerGroup.Consume(ctx, []string{"message_topic"}, mc); err != nil {
			logx.Errorf("消息消费者错误: %v", err)
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (mc *MessageConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (mc *MessageConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (mc *MessageConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			logx.Errorf("解析消息失败: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		logx.Infof("收到消息: %+v", msg)
		// 处理消息逻辑

		session.MarkMessage(message, "")
	}
	return nil
}
