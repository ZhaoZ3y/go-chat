package logic

import (
	"IM/pkg/model"
	"context"
	"fmt"

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

// 会话详情结构体，用于接收JOIN查询结果
type ConversationDetail struct {
	model.Conversations
	TargetName   string `gorm:"column:target_name"`
	TargetAvatar string `gorm:"column:target_avatar"`
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

	// 批量获取用户和群组信息
	userIDs := make([]int64, 0)
	groupIDs := make([]int64, 0)

	for _, conv := range conversations {
		if conv.Type == 0 { // 私聊
			userIDs = append(userIDs, conv.TargetId)
		} else { // 群聊
			groupIDs = append(groupIDs, conv.TargetId)
		}
	}

	// 批量查询用户信息
	var users []model.User
	userMap := make(map[int64]model.User)
	if len(userIDs) > 0 {
		if err := l.svcCtx.DB.Select("id, nickname, avatar").
			Where("id IN ?", userIDs).
			Find(&users).Error; err != nil {
			l.Logger.Errorf("批量查询用户信息失败: %v", err)
		} else {
			for _, user := range users {
				userMap[user.Id] = user
			}
		}
	}

	// 批量查询群组信息
	var groups []model.Groups
	groupMap := make(map[int64]model.Groups)
	if len(groupIDs) > 0 {
		if err := l.svcCtx.DB.Select("id, name, avatar").
			Where("id IN ?", groupIDs).
			Find(&groups).Error; err != nil {
			l.Logger.Errorf("批量查询群组信息失败: %v", err)
		} else {
			for _, group := range groups {
				groupMap[group.Id] = group
			}
		}
	}

	// 构建返回结果
	chatConversations := make([]*chat.Conversation, 0, len(conversations))

	for _, conv := range conversations {
		targetName := ""
		targetAvatar := ""

		if conv.Type == 0 { // 私聊
			if user, exists := userMap[conv.TargetId]; exists {
				targetName = user.Nickname
				targetAvatar = user.Avatar
			} else {
				targetName = fmt.Sprintf("用户%d", conv.TargetId)
			}
		} else { // 群聊
			if group, exists := groupMap[conv.TargetId]; exists {
				targetName = group.Name
				targetAvatar = group.Avatar
			} else {
				targetName = fmt.Sprintf("群组%d", conv.TargetId)
			}
		}

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
			TargetName:      targetName,
			TargetAvatar:    targetAvatar,
		})
	}

	return &chat.GetConversationListResponse{
		Conversations: chatConversations,
	}, nil
}
