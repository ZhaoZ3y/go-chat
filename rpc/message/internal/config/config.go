package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf
	RocketMQ   RocketMQConf
}

type RocketMQConf struct {
	NameSrvAddrs  []string
	ProducerGroup string
	ConsumerGroup string
}
