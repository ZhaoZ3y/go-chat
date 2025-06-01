package main

import (
	"IM/api/router"
	"IM/pkg/websocket"
)

func main() {
	// 初始化websocket Hub
	hub := websocket.NewHub()
	go hub.Run() // 启动WebSocket Hub

	// 初始化路由
	api := router.SetRouter()
	api.Run(":8080") // 启动API服务，监听8080端口
}
