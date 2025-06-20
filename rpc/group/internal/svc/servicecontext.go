package svc

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"IM/pkg/mq/notify"
	"IM/rpc/group/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	NotifyService notify.NotifyService
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移数据表
	db.AutoMigrate(&model.Friends{}, &model.FriendRequests{})

	// 初始化Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.CustomRedis.Host, c.CustomRedis.Port),
		Password: c.CustomRedis.Password,
		DB:       c.CustomRedis.DB,
	})

	// 初始化Kafka客户端
	kafkaClient, err := mq.NewKafkaClient([]string{"kafka:9092"})
	if err != nil {
		panic("failed to connect kafka")
	}

	// 初始化通知服务
	notifyService := notify.NewNotifyService(kafkaClient)

	return &ServiceContext{
		Config:        c,
		DB:            db,
		Redis:         rdb,
		NotifyService: notifyService,
	}
}
