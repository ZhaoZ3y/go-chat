package mq

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type KafkaClient struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	Brokers  []string
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

func NewKafkaClient(brokers []string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	// 创建生产者
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// 创建消费者
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &KafkaClient{
		Producer: producer,
		Consumer: consumer,
		Brokers:  brokers,
	}, nil
}

func (k *KafkaClient) SendMessage(topic string, event *MessageEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 设置分区键
	var partitionKey string
	if event.ChatType == 0 { // 私聊
		partitionKey = string(rune(event.ToUserID))
	} else { // 群聊
		partitionKey = string(rune(event.GroupID))
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(partitionKey),
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("type"),
				Value: []byte(event.Type),
			},
		},
	}

	_, _, err = k.Producer.SendMessage(msg)
	return err
}

func (k *KafkaClient) CreateConsumer(topic string, consumerGroup string, handler func(*MessageEvent) error) error {
	partitionConsumer, err := k.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		defer partitionConsumer.Close()
		for {
			select {
			case message := <-partitionConsumer.Messages():
				var event MessageEvent
				if err := json.Unmarshal(message.Value, &event); err != nil {
					log.Printf("反序列化消息失败: %v", err)
					continue
				}

				if err := handler(&event); err != nil {
					log.Printf("处理消息失败: %v", err)
				}
			case err := <-partitionConsumer.Errors():
				log.Printf("消费消息错误: %v", err)
			}
		}
	}()

	return nil
}

func (k *KafkaClient) Close() error {
	if k.Producer != nil {
		k.Producer.Close()
	}
	if k.Consumer != nil {
		k.Consumer.Close()
	}
	return nil
}
