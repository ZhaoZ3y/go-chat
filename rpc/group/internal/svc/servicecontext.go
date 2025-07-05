package svc

import (
	"IM/pkg/model"
	"IM/pkg/utils/scheduler"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/config"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	MuteScheduler *scheduler.MuteScheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移数据表
	db.AutoMigrate(&model.Groups{}, &model.GroupMembers{}, &model.JoinGroupApplications{}, model.GroupNotification{})

	// 初始化Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.CustomRedis.Host, c.CustomRedis.Port),
		Password: c.CustomRedis.Password,
		DB:       c.CustomRedis.DB,
	})

	svc := &ServiceContext{
		Config:        c,
		DB:            db,
		Redis:         rdb,
		MuteScheduler: scheduler.NewMuteScheduler(),
	}

	svc.initMuteSchedulers()
	return svc
}

// initMuteSchedulers 初始化禁言调度器
func (svc *ServiceContext) initMuteSchedulers() {
	var muted []model.GroupMembers
	svc.DB.Where("status = ?", group.MemberStatus_MEMBER_STATUS_MUTED).Find(&muted)
	for _, m := range muted {
		key := fmt.Sprintf("group:mute:%d:%d", m.GroupId, m.UserId)
		val, err := svc.Redis.Get(context.Background(), key).Result()
		if err != nil || val == "" {
			svc.DB.Model(&model.GroupMembers{}).
				Where("group_id = ? AND user_id = ?", m.GroupId, m.UserId).
				Update("status", group.MemberStatus_MEMBER_STATUS_NORMAL)
			continue
		}
		until, _ := strconv.ParseInt(val, 10, 64)
		delay := time.Until(time.Unix(until, 0))
		if delay > 0 {
			groupId, userId := m.GroupId, m.UserId
			svc.MuteScheduler.Register(groupId, userId, delay, func() {
				SyncUnmuteStatus(svc, groupId, userId)
			})
		}
	}
}

// SyncUnmuteStatus 检查禁言状态并自动解禁
func SyncUnmuteStatus(svc *ServiceContext, groupId, userId int64) {
	key := fmt.Sprintf("group:mute:%d:%d", groupId, userId)
	exists, err := svc.Redis.Exists(context.Background(), key).Result()
	if err != nil || exists == 1 {
		return
	}
	svc.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND user_id = ?", groupId, userId).
		Update("status", group.MemberStatus_MEMBER_STATUS_NORMAL)
}
