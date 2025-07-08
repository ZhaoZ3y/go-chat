package db

import (
	"IM/pkg/config"
	"IM/pkg/mq"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(c config.Config) (*gorm.DB, *redis.Client, *mq.KafkaClient, error) {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	// 初始化Kafka
	kafkaClient, err := mq.NewKafkaClient(c.Kafka.Brokers)
	if err != nil {
		panic("failed to connect kafka: " + err.Error())
	}

	return db, rdb, kafkaClient, err
}
