package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Redis      RedisConfig
	Salt       string // 密码加盐
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
