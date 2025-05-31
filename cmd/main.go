package main

import "IM/api/router"

func main() {
	api := router.SetRouter()
	api.Run(":8080") // 启动API服务，监听8080端口
}
