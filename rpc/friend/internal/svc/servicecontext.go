package svc

import (
	"IM/pkg/model"
	"IM/rpc/friend/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
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
		Addr:     c.Redis.Host + ":" + c.Redis.Port,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rdb,
	}
}
