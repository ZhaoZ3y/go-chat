package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DataSource  string
	CustomRedis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	MinIO struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		BucketName      string
	}
	FileStorage struct {
		BaseURL string // 文件存储的基础URL，通常是CDN或反向代理地址

	}
}
