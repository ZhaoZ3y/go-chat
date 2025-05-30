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

type MarkMessageReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkMessageReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkMessageReadLogic {
	return &MarkMessageReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记消息已读
func (l *MarkMessageReadLogic) MarkMessageRead(in *chat.MarkMessageReadRequest) (*chat.MarkMessageReadResponse, error) {
	// 清零未读消息数
	err := l.svcCtx.DB.Model(&model.Conversations{}).
		Where("user_id = ? AND target_id = ? AND type = ?", in.UserId, in.TargetId, int8(in.ChatType)).
		Update("unread_count", 0).Error

	if err != nil {
		return &chat.MarkMessageReadResponse{
			Success: false,
			Message: "标记已读失败",
		}, nil
	}

	// 发送已读事件到RocketMQ
	event := &mq.MessageEvent{
		Type:       mq.EventMessageRead,
		FromUserID: in.UserId,
		ToUserID:   in.TargetId,
		ChatType:   int8(in.ChatType),
		CreateAt:   time.Now().Unix(),
		Data: map[string]interface{}{
			"last_read_message_id": in.LastReadMessageId,
		},
	}

	if err := l.svcCtx.RocketMQ.SendMessage(mq.TopicNotify, event); err != nil {
		l.Logger.Errorf("发送已读事件失败: %v", err)
	}

	return &chat.MarkMessageReadResponse{
		Success: true,
		Message: "标记已读成功",
	}, nil
}
