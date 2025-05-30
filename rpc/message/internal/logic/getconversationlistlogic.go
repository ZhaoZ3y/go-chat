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
	var total int64

	// 统计总数
	if err := l.svcCtx.DB.Model(&model.Conversations{}).Where("user_id = ?", in.UserId).Count(&total).Error; err != nil {
		l.Logger.Errorf("统计会话总数失败: %v", err)
		return &chat.GetConversationListResponse{}, err
	}

	// 查询会话列表
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).
		Order("last_message_time DESC").
		Find(&conversations).Error

	if err != nil {
		l.Logger.Errorf("查询会话列表失败: %v", err)
		return &chat.GetConversationListResponse{}, err
	}

	// 转换为proto格式
	var chatConversations []*chat.Conversation
	for _, conv := range conversations {
		chatConversations = append(chatConversations, &chat.Conversation{
			Id:              conv.Id,
			UserId:          conv.UserId,
			TargetId:        conv.TargetId,
			Type:            chat.ChatType(conv.Type),
			LastMessage:     conv.LastMessage,
			LastMessageTime: conv.LastMessageTime,
			UnreadCount:     int32(conv.UnreadCount),
			CreateAt:        conv.CreateAt,
			UpdateAt:        conv.UpdateAt,
		})
	}

	return &chat.GetConversationListResponse{
		Conversations: chatConversations,
		Total:         total,
	}, nil
}
