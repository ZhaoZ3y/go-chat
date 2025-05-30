package mq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketMQClient struct {
	Producer     rocketmq.Producer
	Consumer     rocketmq.PushConsumer
	NameSrvAddrs []string
}

type MessageEvent struct {
	Type        string      `json:"type"`
	MessageID   int64       `json:"message_id"`
	FromUserID  int64       `json:"from_user_id"`
	ToUserID    int64       `json:"to_user_id"`
	GroupID     int64       `json:"group_id"`
	ChatType    int8        `json:"chat_type"`
	Content     string      `json:"content"`
	MessageType int8        `json:"message_type"`
	Extra       string      `json:"extra"`
	CreateAt    int64       `json:"create_at"`
	Data        interface{} `json:"data,omitempty"`
}

const (
	TopicMessage = "IM_MESSAGE"
	TopicNotify  = "IM_NOTIFY"

	EventNewMessage    = "new_message"
	EventMessageRead   = "message_read"
	EventMessageRecall = "message_recall"
	EventMessageDelete = "message_delete"
)

func NewRocketMQClient(nameSrvAddrs []string) (*RocketMQClient, error) {
	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nameSrvAddrs),
		producer.WithRetry(2),
		producer.WithGroupName("im_producer_group"),
	)
	if err != nil {
		return nil, err
	}

	if err := p.Start(); err != nil {
		return nil, err
	}

	return &RocketMQClient{
		Producer:     p,
		NameSrvAddrs: nameSrvAddrs,
	}, nil
}

func (r *RocketMQClient) SendMessage(topic string, event *MessageEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &primitive.Message{
		Topic: topic,
		Body:  data,
	}

	// 设置消息标签
	msg.WithTag(event.Type)

	// 私聊消息按用户ID分区，群聊消息按群ID分区
	if event.ChatType == 0 { // 私聊
		msg.WithShardingKey(string(rune(event.ToUserID)))
	} else { // 群聊
		msg.WithShardingKey(string(rune(event.GroupID)))
	}

	_, err = r.Producer.SendSync(context.Background(), msg)
	return err
}

func (r *RocketMQClient) CreateConsumer(groupName string, handler func(*MessageEvent) error) error {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(r.NameSrvAddrs),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(groupName),
	)
	if err != nil {
		return err
	}

	err = c.Subscribe(TopicMessage, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var event MessageEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("反序列化消息失败: %v", err)
				continue
			}

			if err := handler(&event); err != nil {
				log.Printf("处理消息失败: %v", err)
				return consumer.ConsumeRetryLater, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		return err
	}

	r.Consumer = c
	return c.Start()
}

func (r *RocketMQClient) Close() error {
	if r.Producer != nil {
		r.Producer.Shutdown()
	}
	if r.Consumer != nil {
		r.Consumer.Shutdown()
	}
	return nil
}
