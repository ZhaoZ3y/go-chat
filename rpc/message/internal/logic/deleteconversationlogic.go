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
	result := l.svcCtx.DB.
		Where("user_id = ? AND target_id = ? AND type = ?", in.UserId, in.TargetId, int8(in.ChatType)).
		Delete(&model.Conversations{})

	if result.Error != nil {
		l.Logger.Errorf("软删除会话失败: user_id=%d, target_id=%d, error=%v", in.UserId, in.TargetId, result.Error)
		return &chat.DeleteConversationResponse{
			Success: false,
			Message: "删除会话失败",
		}, result.Error // 返回原始错误以便调试
	}

	return &chat.DeleteConversationResponse{
		Success: true,
		Message: "会话已删除",
	}, nil
}
