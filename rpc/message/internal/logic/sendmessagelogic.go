package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	_const "IM/pkg/utils/const"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	// --- 参数校验 (保持不变) ---
	if in.Content == "" {
		return &chat.SendMessageResponse{Success: false, Message: "消息内容不能为空"}, nil
	}
	if in.ChatType == 0 && in.ToUserId == 0 {
		return &chat.SendMessageResponse{Success: false, Message: "私聊必须指定接收者ID"}, nil
	}
	if in.ChatType == 1 && in.GroupId == 0 {
		return &chat.SendMessageResponse{Success: false, Message: "群聊必须指定群组ID"}, nil
	}

	// --- 群聊权限检查 (保持不变) ---
	if in.ChatType == 1 {
		var member model.GroupMembers
		if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", in.GroupId, in.FromUserId).First(&member).Error; err != nil {
			return &chat.SendMessageResponse{Success: false, Message: "您不是该群组成员或已被禁言"}, nil
		}
	}

	// 准备消息实体
	message := &model.Messages{
		FromUserId: in.FromUserId,
		ToUserId:   in.ToUserId,
		GroupId:    in.GroupId,
		Type:       int64(in.Type),
		Content:    in.Content,
		Extra:      in.Extra,
		ChatType:   int64(in.ChatType),
		Status:     _const.MsgNormal,
	}

	// 事务前打印发送内容
	l.Logger.Infof("发送消息请求: from=%d to=%d group=%d chatType=%d content=%s",
		in.FromUserId, in.ToUserId, in.GroupId, in.ChatType, in.Content)

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建消息记录
		// 这会通过 GORM 的钩子自动填充 CreateAt
		if err := tx.Create(message).Error; err != nil {
			return fmt.Errorf("保存消息到数据库失败: %w", err)
		}

		// 2. 更新或创建相关会话（私聊/群聊）
		// 调用重构后的内部方法
		if err := l.updateOrCreateConversation(tx, message); err != nil {
			return fmt.Errorf("更新会话失败: %w", err)
		}

		if message.ChatType == 0 {
			// 私聊：写入接收方（ToUserId）
			if err := tx.Create(&model.MessageUserStates{
				MessageId: message.Id,
				UserId:    message.ToUserId,
				IsDeleted: false,
			}).Error; err != nil {
				return fmt.Errorf("写入 message_user_states (to_user) 失败: %w", err)
			}

			// 可选：如果你希望发送方也能有记录，可加上
			if err := tx.Create(&model.MessageUserStates{
				MessageId: message.Id,
				UserId:    message.FromUserId,
				IsDeleted: false,
			}).Error; err != nil {
				return fmt.Errorf("写入 message_user_states (from_user) 失败: %w", err)
			}

		} else {
			// 群聊：写入所有群成员
			var memberIDs []int64
			if err := tx.Model(&model.GroupMembers{}).
				Where("group_id = ? AND status = 1", message.GroupId).
				Pluck("user_id", &memberIDs).Error; err != nil {
				return fmt.Errorf("写入群聊 message_user_states 前拉取成员失败: %w", err)
			}

			var stateRecords []model.MessageUserStates
			for _, uid := range memberIDs {
				stateRecords = append(stateRecords, model.MessageUserStates{
					MessageId: message.Id,
					UserId:    uid,
					IsDeleted: false,
				})
			}

			if len(stateRecords) > 0 {
				if err := tx.Create(&stateRecords).Error; err != nil {
					return fmt.Errorf("批量写入 message_user_states 失败: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		l.Logger.Errorf("发送消息事务失败: %v", err)
		return &chat.SendMessageResponse{
			Success: false,
			Message: "发送消息失败",
		}, nil
	}

	event := &mq.MessageEvent{
		Type:        mq.EventNewMessage,
		MessageID:   message.Id,
		FromUserID:  message.FromUserId,
		ToUserID:    message.ToUserId,
		GroupID:     message.GroupId,
		ChatType:    message.ChatType,
		Content:     message.Content,
		MessageType: message.Type,
		Extra:       message.Extra,
		CreateAt:    message.CreateAt,
	}

	if err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, event); err != nil {
		l.Logger.Errorf("发送消息事件到 Kafka 失败 (但消息已存库): %v", err)
	}

	// Kafka 事件发出后打印
	l.Logger.Infof("Kafka事件已发送: messageID=%d", message.Id)

	return &chat.SendMessageResponse{
		MessageId: message.Id,
		Success:   true,
		Message:   "发送成功",
	}, nil
}

// 更新会话记录
func (l *SendMessageLogic) updateOrCreateConversation(tx *gorm.DB, message *model.Messages) error {
	now := message.CreateAt // 统一使用消息创建时的时间戳

	lastMessageContent := message.Content
	if message.Type == _const.MsgTypeImage {
		lastMessageContent = "[图片]"
	} else if message.Type == _const.MsgTypeAudio {
		lastMessageContent = "[语音]"
	} else if message.Type == _const.MsgTypeFile {
		lastMessageContent = "[文件]"
	} else if message.Type == _const.MsgTypeVideo {
		lastMessageContent = "[视频]"
	}

	if message.ChatType == 0 { // 私聊
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "target_id"}, {Name: "type"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_message_id", "last_message_time", "last_message"}),
		}).Create(&model.Conversations{
			UserId:          message.FromUserId,
			TargetId:        message.ToUserId,
			Type:            0,
			LastMessageID:   message.Id,
			LastMessageTime: now,
			LastMessage:     lastMessageContent,
		}).Error
		if err != nil {
			return fmt.Errorf("更新发送方会话失败: %w", err)
		}

		// 更新或创建接收方的会话（未读数+1）
		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "target_id"}, {Name: "type"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"last_message_id":   message.Id,
				"last_message_time": now,
				"last_message":      lastMessageContent,
				"unread_count":      gorm.Expr("unread_count + 1"),
			}),
		}).Create(&model.Conversations{
			UserId:          message.ToUserId,
			TargetId:        message.FromUserId,
			Type:            0,
			LastMessageID:   message.Id,
			LastMessageTime: now,
			LastMessage:     lastMessageContent,
			UnreadCount:     1,
		}).Error
		if err != nil {
			return fmt.Errorf("更新接收方会话失败: %w", err)
		}

	} else { // 群聊
		var memberIDs []int64
		err := tx.Model(&model.GroupMembers{}).
			Where("group_id = ? AND status = 1", message.GroupId).
			Pluck("user_id", &memberIDs).Error
		if err != nil {
			return fmt.Errorf("获取群成员失败: %w", err)
		}

		for _, memberID := range memberIDs {
			// 对非发送者，未读数+1
			unreadIncrement := gorm.Expr("unread_count + 1")
			initialUnread := 1
			if memberID == message.FromUserId {
				// 发送者自己的会话，未读数不变
				unreadIncrement = gorm.Expr("unread_count")
				initialUnread = 0
			}

			err = tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}, {Name: "target_id"}, {Name: "type"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"last_message_id":   message.Id,
					"last_message_time": now,
					"last_message":      lastMessageContent,
					"unread_count":      unreadIncrement,
				}),
			}).Create(&model.Conversations{
				UserId:          memberID,
				TargetId:        message.GroupId,
				Type:            1,
				LastMessageID:   message.Id,
				LastMessageTime: now,
				LastMessage:     lastMessageContent,
				UnreadCount:     initialUnread,
			}).Error

			if err != nil {
				l.Logger.Errorf("更新群成员 %d 的会话失败: %v", memberID, err)
			}
		}
	}
	return nil
}
