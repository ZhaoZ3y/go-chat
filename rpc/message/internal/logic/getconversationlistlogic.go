package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationListLogic {
	return &GetConversationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话列表
func (l *GetConversationListLogic) GetConversationList(in *chat.GetConversationListRequest) (*chat.GetConversationListResponse, error) {
	var conversations []model.Conversations
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).
		Order("is_pinned DESC, update_at DESC").
		Find(&conversations).Error

	if err != nil {
		l.Logger.Errorf("查询用户 %d 的会话列表失败: %v", in.UserId, err)
		return nil, err
	}

	chatConversations := make([]*chat.Conversation, 0, len(conversations))
	for _, conv := range conversations {
		chatConversations = append(chatConversations, &chat.Conversation{
			Id:              conv.Id,
			UserId:          conv.UserId,
			TargetId:        conv.TargetId,
			Type:            chat.ChatType(conv.Type),
			LastMessageId:   conv.LastMessageID,
			LastMessage:     conv.LastMessage,
			LastMessageTime: conv.LastMessageTime,
			UnreadCount:     int32(conv.UnreadCount),
			IsPinned:        conv.IsPinned,
			CreateAt:        conv.CreateAt,
			UpdateAt:        conv.UpdateAt,
		})
	}

	return &chat.GetConversationListResponse{
		Conversations: chatConversations,
	}, nil
}
