package svc

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"IM/rpc/message/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Kafka  *mq.KafkaClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移表结构
	db.AutoMigrate(&model.Messages{}, &model.Conversations{})

	// 初始化Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.CustomRedis.Host, c.CustomRedis.Port),
		Password: c.CustomRedis.Password,
		DB:       c.CustomRedis.DB,
	})

	// 初始化Kafka
	kafkaClient, err := mq.NewKafkaClient(c.Kafka.Brokers)
	if err != nil {
		panic("failed to connect kafka: " + err.Error())
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rdb,
		Kafka:  kafkaClient,
	}
}
