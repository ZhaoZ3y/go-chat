package logic

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"

	"IM/rpc/notify/internal/svc"
	"IM/rpc/notify/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNotificationLogic {
	return &SendNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送通知接口，将通知消息推送到 Kafka 等消息队列
func (l *SendNotificationLogic) SendNotification(in *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	// 设置时间戳
	if in.Notification.Timestamp == 0 {
		in.Notification.Timestamp = time.Now().Unix()
	}

	// 序列化通知消息
	data, err := json.Marshal(in.Notification)
	if err != nil {
		l.Errorf("序列化通知消息失败: %v", err)
		return &notification.SendNotificationResponse{
			Success:  false,
			ErrorMsg: "序列化错误",
		}, nil
	}

	// 发送到Kafka
	msg := &sarama.ProducerMessage{
		Topic:     l.svcCtx.Config.Kafka.ProducerTopic,
		Key:       sarama.StringEncoder(string(rune(in.Notification.UserId))),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	}

	partition, offset, err := l.svcCtx.KafkaProducer.SendMessage(msg)
	if err != nil {
		l.Errorf("发送Kafka消息失败: %v", err)
		return &notification.SendNotificationResponse{
			Success:  false,
			ErrorMsg: "消息队列发送失败",
		}, nil
	}

	l.Infof("通知消息发送成功: UserID=%d, Type=%v, Partition=%d, Offset=%d",
		in.Notification.UserId, in.Notification.Type, partition, offset)

	return &notification.SendNotificationResponse{
		Success: true,
	}, nil
}
