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
	if in.Content == "" {
		return &chat.SendMessageResponse{Success: false, Message: "消息内容不能为空"}, nil
	}
	if in.ChatType == 0 && in.ToUserId == 0 {
		return &chat.SendMessageResponse{Success: false, Message: "私聊必须指定接收者ID"}, nil
	}
	if in.ChatType == 1 && in.GroupId == 0 {
		return &chat.SendMessageResponse{Success: false, Message: "群聊必须指定群组ID"}, nil
	}

	// 检查是否为群成员
	if in.ChatType == 1 {
		var member model.GroupMembers
		if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.FromUserId).First(&member).Error; err != nil {
			return &chat.SendMessageResponse{Success: false, Message: "您不是该群组成员"}, nil
		}
	}

	// 创建消息
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

	// 打日志
	l.Logger.Infof("发送消息请求: from=%d to=%d group=%d chatType=%d content=%s",
		in.FromUserId, in.ToUserId, in.GroupId, in.ChatType, in.Content)

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 保存消息（自动生成 ID）
		if err := tx.Create(message).Error; err != nil {
			return fmt.Errorf("保存消息失败: %w", err)
		}

		// 更新会话
		if err := l.updateOrCreateConversation(tx, message); err != nil {
			return fmt.Errorf("更新会话失败: %w", err)
		}

		// 插入 message_user_states，防止重复插入
		if message.ChatType == 0 {
			// 私聊：from 和 to 各插一条
			users := []int64{message.FromUserId, message.ToUserId}
			for _, uid := range users {
				state := model.MessageUserStates{
					MessageId: message.Id,
					UserId:    uid,
					IsDeleted: false,
				}
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}},
					DoNothing: true,
				}).Create(&state).Error; err != nil {
					return fmt.Errorf("写入 message_user_states 失败: %w", err)
				}
			}
		} else {
			// 群聊：群成员都插
			var memberIDs []int64
			if err := tx.Model(&model.GroupMembers{}).
				Where("group_id = ? AND status = 1", message.GroupId).
				Pluck("user_id", &memberIDs).Error; err != nil {
				return fmt.Errorf("拉取群成员失败: %w", err)
			}

			var states []model.MessageUserStates
			for _, uid := range memberIDs {
				states = append(states, model.MessageUserStates{
					MessageId: message.Id,
					UserId:    uid,
					IsDeleted: false,
				})
			}

			if len(states) > 0 {
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}},
					DoNothing: true,
				}).Create(&states).Error; err != nil {
					return fmt.Errorf("批量写入 message_user_states 失败: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		l.Logger.Errorf("发送消息事务失败: %v", err)
		return &chat.SendMessageResponse{Success: false, Message: "发送失败"}, nil
	}

	// 推送 Kafka 消息
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
		l.Logger.Errorf("Kafka 推送失败（但消息已存库）: %v", err)
	}

	l.Logger.Infof("消息发送成功，Kafka事件已发出: messageID=%d", message.Id)
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
		// 添加调试日志
		l.Logger.Infof("开始处理群聊会话: groupId=%d, messageId=%d", message.GroupId, message.Id)

		var memberIDs []int64
		err := tx.Model(&model.GroupMembers{}).
			Where("group_id = ?", message.GroupId).
			Pluck("user_id", &memberIDs).Error
		if err != nil {
			l.Logger.Errorf("获取群成员失败: groupId=%d, err=%v", message.GroupId, err)
			return fmt.Errorf("获取群成员失败: %w", err)
		}

		l.Logger.Infof("找到群成员: groupId=%d, members=%v", message.GroupId, memberIDs)

		if len(memberIDs) == 0 {
			l.Logger.Errorf("群组 %d 没有找到活跃成员", message.GroupId)
			return nil
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

			l.Logger.Infof("处理群成员会话: userId=%d, groupId=%d, isFromUser=%v",
				memberID, message.GroupId, memberID == message.FromUserId)

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
				return fmt.Errorf("更新群成员 %d 的会话失败: %w", memberID, err)
			} else {
				l.Logger.Infof("成功更新群成员会话: userId=%d, groupId=%d", memberID, message.GroupId)
			}
		}

		l.Logger.Infof("群聊会话处理完成: groupId=%d, 处理成员数=%d", message.GroupId, len(memberIDs))
	}
	return nil
}
