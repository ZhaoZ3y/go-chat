package logic

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"

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
	data, err := json.Marshal(in.Notification)
	if err != nil {
		l.Errorf("failed to marshal notification: %v", err)
		return &notification.SendNotificationResponse{
			Success:  false,
			ErrorMsg: "marshal error",
		}, nil
	}

	msg := &sarama.ProducerMessage{
		Topic: l.svcCtx.Config.Kafka.ProducerTopic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := l.svcCtx.KafkaProducer.SendMessage(msg)
	if err != nil {
		l.Errorf("failed to send kafka message: %v", err)
		return &notification.SendNotificationResponse{
			Success:  false,
			ErrorMsg: "kafka send error",
		}, nil
	}

	l.Infof("sent kafka message partition=%d offset=%d", partition, offset)
	return &notification.SendNotificationResponse{
		Success: true,
	}, nil
}
