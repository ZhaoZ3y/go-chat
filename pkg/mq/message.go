package mq

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	TopicMessage      = "im_message"
	TopicNotification = "im_notification"
	TopicUserStatus   = "im_user_status"
)

const (
	EventNewMessage    = "new_message"
	EventUserOnline    = "user_online"
	EventUserOffline   = "user_offline"
	EventHeartbeat     = "heartbeat"
	EventMessageRead   = "message_read"
	EventNotification  = "notification"
	EventMessageRecall = "message_recall"
)

type MessageEvent struct {
	Type              string `json:"type"`
	MessageID         int64  `json:"message_id"`
	FromUserID        int64  `json:"from_user_id"`
	ToUserID          int64  `json:"to_user_id"`
	GroupID           int64  `json:"group_id"`
	ChatType          int64  `json:"chat_type"`
	Content           string `json:"content"`
	MessageType       int64  `json:"message_type"`
	Extra             string `json:"extra"`
	CreateAt          int64  `json:"create_at"`
	LastReadMessageID int64  `json:"last_read_message_id,omitempty"`
}

type UserStatusEvent struct {
	Type      string `json:"type"`
	UserID    int64  `json:"user_id"`
	Status    int8   `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type RichMessageEvent struct {
	MessageEvent         // 嵌入基础事件，JSON序列化时会将其字段平铺
	Recipients   []int64 `json:"recipients"` // 所有需要接收此事件的用户ID列表
}

type NotificationEvent struct {
	Type     string `json:"type"`
	UserID   int64  `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Extra    string `json:"extra"`
	CreateAt int64  `json:"create_at"`
}

type KafkaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewKafkaClient(brokers []string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kafka 生产者失败: %w", err)
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 Kafka 消费者失败: %w", err)
	}

	return &KafkaClient{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (kc *KafkaClient) Close() {
	if kc.producer != nil {
		kc.producer.Close()
	}
	if kc.consumer != nil {
		kc.consumer.Close()
	}
}

func (kc *KafkaClient) SendMessage(topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("序列化事件失败: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := kc.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	logx.Infof("消息发送成功 topic:%s partition:%d offset:%d", topic, partition, offset)
	return nil
}

func (kc *KafkaClient) ConsumeMessages(topic string, handler func([]byte) error) error {
	partitions, err := kc.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("获取分区失败: %w", err)
	}

	for _, partition := range partitions {
		pc, err := kc.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("消费分区失败: %w", err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for message := range pc.Messages() {
				if err := handler(message.Value); err != nil {
					logx.Errorf("处理消息失败: %v", err)
				}
			}
		}(pc)
	}

	return nil
}
