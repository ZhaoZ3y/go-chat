package logic

import (
	"IM/pkg/model"
	"context"
	"time"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送消息
func (l *SendMessageLogic) SendMessage(in *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	// 创建消息记录
	message := &model.Messages{
		FromUserId: in.FromUserId,
		ToUserId:   in.ToUserId,
		GroupId:    in.GroupId,
		Type:       int8(in.Type),
		Content:    in.Content,
		Extra:      in.Extra,
		ChatType:   int8(in.ChatType),
		Status:     1,
		CreateAt:   time.Now().Unix(),
	}

	// 保存消息到数据库
	if err := l.svcCtx.DB.Create(message).Error; err != nil {
		l.Logger.Errorf("保存消息失败: %v", err)
		return &chat.SendMessageResponse{
			Success: false,
			Message: "发送消息失败",
		}, nil
	}

	// 更新或创建会话记录
	if err := l.updateConversation(in, message.Content); err != nil {
		l.Logger.Errorf("更新会话失败: %v", err)
	}

	return &chat.SendMessageResponse{
		MessageId: message.Id,
		Success:   true,
		Message:   "发送成功",
	}, nil
}

func (l *SendMessageLogic) updateConversation(req *chat.SendMessageRequest, content string) error {
	now := time.Now().Unix()

	if req.ChatType == chat.ChatType_PRIVATE {
		// 私聊：更新发送者和接收者的会话记录
		conversations := []model.Conversations{
			{
				UserId:          req.FromUserId,
				TargetId:        req.ToUserId,
				Type:            int8(chat.ChatType_PRIVATE),
				LastMessage:     content,
				LastMessageTime: now,
			},
			{
				UserId:          req.ToUserId,
				TargetId:        req.FromUserId,
				Type:            int8(chat.ChatType_PRIVATE),
				LastMessage:     content,
				LastMessageTime: now,
				UnreadCount:     1, // 接收者未读消息+1
			},
		}

		for _, conv := range conversations {
			l.svcCtx.DB.Where("user_id = ? AND target_id = ? AND type = ?",
				conv.UserId, conv.TargetId, conv.Type).
				Assign(map[string]interface{}{
					"last_message":      conv.LastMessage,
					"last_message_time": conv.LastMessageTime,
					"unread_count":      conv.UnreadCount,
					"update_at":         now,
				}).FirstOrCreate(&conv)
		}
	} else {
		// 群聊：更新群组会话记录（简化处理，实际应该更新所有群成员的会话）
		conv := model.Conversations{
			UserId:          req.FromUserId,
			TargetId:        req.GroupId,
			Type:            int8(chat.ChatType_GROUP),
			LastMessage:     content,
			LastMessageTime: now,
		}

		l.svcCtx.DB.Where("user_id = ? AND target_id = ? AND type = ?",
			conv.UserId, conv.TargetId, conv.Type).
			Assign(map[string]interface{}{
				"last_message":      conv.LastMessage,
				"last_message_time": conv.LastMessageTime,
				"update_at":         now,
			}).FirstOrCreate(&conv)
	}

	return nil
}
