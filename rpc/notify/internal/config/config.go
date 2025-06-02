package config

import "github.com/zeromicro/go-zero/zrpc"

type KafkaConfig struct {
	Brokers       []string
	ProducerTopic string
	ConsumerGroup string
}

type Config struct {
	zrpc.RpcServerConf
	Kafka KafkaConfig
}
