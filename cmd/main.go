package main

import (
	"IM/api/router"
	"IM/pkg/config"
	"IM/pkg/utils/db"
	"IM/pkg/websocket"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "f", "pkg/config/config.yaml", "the config file")
	flag.Parse()

	var c config.Config
	conf.MustLoad(configFile, &c)

	DB, RedisClient, kafka, err := db.InitDB(c)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	servernode := "server1" // 假设这是当前服务器的节点标识
	hub := websocket.NewHub(DB, kafka, RedisClient, servernode)
	go hub.Run()
	api := router.SetRouter(hub)

	if err := api.Run(":8080"); err != nil {
		log.Fatal("API服务启动失败:", err)
	}
}
