package svc

import (
	"IM/rpc/notify/internal/config"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config        config.Config
	KafkaProducer sarama.SyncProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	producer, err := sarama.NewSyncProducer(c.Kafka.Brokers, nil)
	if err != nil {
		logx.Errorf("failed to create kafka producer: %v", err)
	}
	return &ServiceContext{
		Config:        c,
		KafkaProducer: producer,
	}
}

func (s *ServiceContext) Close() {
	if s.KafkaProducer != nil {
		_ = s.KafkaProducer.Close()
	}
}
