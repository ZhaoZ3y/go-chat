package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPinConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinConversationLogic {
	return &PinConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设置会话置顶状态
func (l *PinConversationLogic) PinConversation(in *chat.PinConversationRequest) (*chat.PinConversationResponse, error) {
	result := l.svcCtx.DB.Model(&model.Conversations{}).
		Where("user_id = ? AND target_id = ? AND type = ?", in.UserId, in.TargetId, int8(in.ChatType)).
		Update("is_pinned", in.IsPinned)

	if result.Error != nil {
		l.Logger.Errorf("为用户 %d 设置会话(target:%d)置顶状态失败: %v", in.UserId, in.TargetId, result.Error)
		return &chat.PinConversationResponse{Success: false, Message: "数据库操作失败"}, result.Error
	}

	if result.RowsAffected == 0 {
		l.Logger.Errorf("用户 %d 尝试设置一个不存在的会话(target:%d)的置顶状态", in.UserId, in.TargetId)
		return &chat.PinConversationResponse{Success: false, Message: "会话不存在，无法置顶"}, nil
	}

	// 因为是单端，所以我们不发送 Kafka 通知。
	// 前端在收到这个成功的响应后，应该主动更新自己的 UI 状态。
	l.Logger.Infof("用户 %d 成功设置会话(target:%d, type:%s)的置顶状态为 %v", in.UserId, in.TargetId, in.ChatType.String(), in.IsPinned)

	return &chat.PinConversationResponse{
		Success: true,
		Message: "操作成功",
	}, nil
}
