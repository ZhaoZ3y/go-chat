package svc

import (
	"IM/pkg/minio"
	"IM/pkg/model"
	"IM/rpc/file/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log" // 使用标准 log 包，因为 go-zero 的 logx 在此阶段可能尚未完全初始化
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	err = db.AutoMigrate(&model.FileRecord{})
	if err != nil {
		log.Fatalf("自动迁移数据库表失败: %v", err)
	}
	log.Println("数据库连接和迁移成功!")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.CustomRedis.Host, c.CustomRedis.Port),
		Password: c.CustomRedis.Password,
		DB:       c.CustomRedis.DB,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}
	log.Println("Redis连接成功!")

	minioClient, err := minio.NewClient(minio.Config{
		Endpoint:        c.MinIO.Endpoint,
		AccessKeyID:     c.MinIO.AccessKeyID,
		SecretAccessKey: c.MinIO.SecretAccessKey,
		UseSSL:          c.MinIO.UseSSL,
		BucketName:      c.MinIO.BucketName,
	})
	if err != nil {
		log.Fatalf("初始化MinIO客户端失败: %v", err)
	}
	log.Println("MinIO客户端初始化成功!")

	lifecycleDays := 7
	err = minioClient.SetBucketLifecycle(context.Background(), lifecycleDays)
	if err != nil {
		log.Printf("警告: 设置MinIO存储桶生命周期策略失败: %v\n", err)
	} else {
		log.Printf("成功设置MinIO存储桶 '%s' 的生命周期策略: %d天后过期\n", c.MinIO.BucketName, lifecycleDays)
	}

	return &ServiceContext{
		Config:      c,
		DB:          db,
		Redis:       rdb,
		MinioClient: minioClient,
	}
}
