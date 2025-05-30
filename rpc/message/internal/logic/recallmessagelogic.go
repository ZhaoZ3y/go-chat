package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"context"
	"time"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecallMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecallMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallMessageLogic {
	return &RecallMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤回消息
func (l *RecallMessageLogic) RecallMessage(in *chat.RecallMessageRequest) (*chat.RecallMessageResponse, error) {
	// 检查消息是否存在且属于当前用户
	var message model.Messages
	err := l.svcCtx.DB.Where("id = ? AND from_user_id = ?", in.MessageId, in.UserId).First(&message).Error
	if err != nil {
		return &chat.RecallMessageResponse{
			Success: false,
			Message: "消息不存在或无权限撤回",
		}, nil
	}

	// 检查撤回时间限制
	if time.Now().Unix()-message.CreateAt > 120 {
		return &chat.RecallMessageResponse{
			Success: false,
			Message: "超过撤回时间限制",
		}, nil
	}

	// 更新消息状态
	err = l.svcCtx.DB.Model(&message).Updates(map[string]interface{}{
		"status":  0,
		"content": "[消息已撤回]",
	}).Error

	if err != nil {
		return &chat.RecallMessageResponse{
			Success: false,
			Message: "撤回消息失败",
		}, nil
	}

	// 发送撤回事件到RocketMQ
	event := &mq.MessageEvent{
		Type:       mq.EventMessageRecall,
		MessageID:  message.Id,
		FromUserID: message.FromUserId,
		ToUserID:   message.ToUserId,
		GroupID:    message.GroupId,
		ChatType:   message.ChatType,
		CreateAt:   time.Now().Unix(),
	}

	if err := l.svcCtx.RocketMQ.SendMessage(mq.TopicNotify, event); err != nil {
		l.Logger.Errorf("发送撤回事件失败: %v", err)
	}

	return &chat.RecallMessageResponse{
		Success: true,
		Message: "撤回成功",
	}, nil
}
