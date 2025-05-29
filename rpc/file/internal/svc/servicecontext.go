package svc

import (
	pkgminio "IM/pkg/minio"
	"IM/pkg/model"
	"IM/rpc/file/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time" // 导入 time 包
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Minio  *pkgminio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("获取通用数据库对象失败: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Files{})
	if err != nil {
		panic("自动迁移 Files 表失败: " + err.Error())
	}

	// 初始化Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	// 初始化MinIO客户端
	minioClient, err := pkgminio.NewMinIOClient(pkgminio.Config{
		Endpoint:        c.MinIO.Endpoint,
		AccessKeyID:     c.MinIO.AccessKeyID,
		SecretAccessKey: c.MinIO.SecretAccessKey,
		UseSSL:          c.MinIO.UseSSL,
		BucketName:      c.MinIO.BucketName,
	})
	if err != nil {
		panic("初始化MinIO客户端失败: " + err.Error())
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rdb,
		Minio:  minioClient,
	}
}
