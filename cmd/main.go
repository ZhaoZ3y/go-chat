package main

import (
	"IM/api/router"
	"IM/pkg/websocket"
	"log"
)

func main() {
	// 初始化WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run() // 启动WebSocket Hub

	// 初始化路由（需要传入hub以便注册WebSocket路由）
	api := router.SetRouter(hub)

	// 启动API服务
	if err := api.Run(":8080"); err != nil {
		log.Fatal("API服务启动失败:", err)
	}
}
