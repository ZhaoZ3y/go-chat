package mq

import (
	"github.com/IBM/sarama"
)

type KafkaClient struct {
	brokers  []string
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewKafkaClient(brokers []string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Version = sarama.V2_1_0_0

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &KafkaClient{
		brokers:  brokers,
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

func (kc *KafkaClient) GetBrokers() []string {
	return kc.brokers
}

func (kc *KafkaClient) GetProducer() sarama.SyncProducer {
	return kc.producer
}

func (kc *KafkaClient) GetConsumer() sarama.Consumer {
	return kc.consumer
}
