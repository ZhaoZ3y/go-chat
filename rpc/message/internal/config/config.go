package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource  string
	CustomRedis RedisConfig
	Kafka       KafkaConf
}

type KafkaConf struct {
	Brokers []string
}
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
