package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationLogic {
	return &DeleteConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除会话
func (l *DeleteConversationLogic) DeleteConversation(in *chat.DeleteConversationRequest) (*chat.DeleteConversationResponse, error) {
	// 执行软删除操作
	err := l.svcCtx.DB.Model(&model.Conversations{}).
		Where("user_id = ? AND target_id = ? AND type = ?", in.UserId, in.TargetId, int8(in.ChatType)).
		Update("deleted", true).Error

	if err != nil {
		l.Logger.Errorf("删除会话失败: %v", err)
		return &chat.DeleteConversationResponse{
			Success: false,
			Message: "删除会话失败",
		}, nil
	}

	return &chat.DeleteConversationResponse{
		Success: true,
		Message: "会话已删除",
	}, nil
}
