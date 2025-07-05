package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DataSource  string
	CustomRedis RedisConfig
	MinIO       MinIOConfig
}

// RedisConfig 对应 YAML 中的 CustomRedis 部分
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// MinIOConfig 对应 YAML 中的 MinIO 部分
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}
