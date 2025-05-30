package svc

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"IM/rpc/message/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	RocketMQ *mq.RocketMQClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移表结构
	db.AutoMigrate(&model.Messages{}, &model.Conversations{})

	// 初始化Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// 初始化RocketMQ
	mqClient, err := mq.NewRocketMQClient(c.RocketMQ.NameSrvAddrs)
	if err != nil {
		panic("failed to connect rocketmq: " + err.Error())
	}

	return &ServiceContext{
		Config:   c,
		DB:       db,
		Redis:    rdb,
		RocketMQ: mqClient,
	}
}
