package main

import (
	"IM/api/router"
	"IM/pkg/config"
	"IM/pkg/database"
	"IM/pkg/mq"
	"IM/pkg/mq/consumer"
	"IM/pkg/websocket"
	"log"
)

func main() {
	var c config.Config
	// 初始化数据库连接
	db, err := database.InitDB(c)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 初始化Kafka客户端
	kafkaClient, err := mq.NewKafkaClient(c.Kafka.Brokers)
	if err != nil {
		log.Fatal("Kafka客户端初始化失败:", err)
	}
	defer kafkaClient.Close()

	// 初始化WebSocket Hub
	hub := websocket.NewHub(db)
	go hub.Run() // 启动WebSocket Hub

	// 初始化WebSocket推送服务
	pushService := websocket.NewPushService(hub)

	// 启动消息消费者
	messageConsumer := consumer.NewMessageConsumer(kafkaClient, pushService)
	go func() {
		if err := messageConsumer.Start(); err != nil {
			log.Printf("消息消费者启动失败: %v", err)
		}
	}()

	// 启动通知消费者
	notifyConsumer := consumer.NewNotifyConsumer(kafkaClient, pushService, db)
	go func() {
		if err := notifyConsumer.Start(); err != nil {
			log.Printf("通知消费者启动失败: %v", err)
		}
	}()

	log.Println("IM服务启动成功")
	log.Println("WebSocket地址: ws://localhost:8080/ws")
	log.Println("API地址: http://localhost:8080")

	// 初始化路由（需要传入hub以便注册WebSocket路由）
	api := router.SetRouter(hub)

	// 启动API服务
	if err := api.Run(":8080"); err != nil {
		log.Fatal("API服务启动失败:", err)
	}
}
