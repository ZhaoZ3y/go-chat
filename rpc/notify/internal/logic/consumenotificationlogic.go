package logic

import (
	"context"

	"IM/rpc/notify/internal/svc"
	"IM/rpc/notify/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumeNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeNotificationLogic {
	return &ConsumeNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 消费通知接口，消费端实现消息接收处理（可选）
func (l *ConsumeNotificationLogic) ConsumeNotification(in *notification.ConsumeNotificationRequest) (*notification.ConsumeNotificationResponse, error) {
	l.Infof("consuming notification: %+v", in.Notification)
	// 这里可以接入 WebSocket 推送、入库等逻辑
	return &notification.ConsumeNotificationResponse{Success: true}, nil
}
