package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	var unreadMessages []model.Messages

	subQuery := l.svcCtx.DB.Model(&model.MessageReadReceipts{}).
		Select("message_id").
		Where("user_id = ?", in.UserId)

	query := l.svcCtx.DB.Model(&model.Messages{}).
		Where("id <= ?", in.LastReadMessageId).
		Where("id NOT IN (?)", subQuery). // 排除已读的
		Where("status = ?", 0)            // 只处理正常状态的消息

	if in.ChatType == chat.ChatType_PRIVATE {
		query = query.Where("to_user_id = ? AND from_user_id = ? AND chat_type = ?", in.UserId, in.TargetId, in.ChatType)
	} else {
		query = query.Where("group_id = ? AND from_user_id != ? AND chat_type = ?", in.TargetId, in.UserId, in.ChatType)
	}

	if err := query.Find(&unreadMessages).Error; err != nil {
		l.Logger.Errorf("标记已读失败：查找未读消息失败: %v", err)
		return nil, err
	}

	if len(unreadMessages) == 0 {
		// 即使没有消息需要标记，也应该清空会话的未读数，以防万一
		_ = l.svcCtx.DB.Model(&model.Conversations{}).
			Where("user_id = ? AND target_id = ? AND type = ?", in.UserId, in.TargetId, in.ChatType).
			Update("unread_count", 0).Error

		l.Logger.Infof("用户 %d 在会话 %d 中没有新的未读消息需要标记", in.UserId, in.TargetId)
		return &chat.MarkMessageReadResponse{Success: true, Message: "没有新的未读消息"}, nil
	}

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		receipts := make([]model.MessageReadReceipts, len(unreadMessages))
		for i, msg := range unreadMessages {
			receipts[i] = model.MessageReadReceipts{
				MessageId: msg.Id,
				UserId:    in.UserId,
				GroupId:   msg.GroupId,
				ChatType:  int8(msg.ChatType),
			}
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&receipts).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Conversations{}).
			Where("user_id = ? AND target_id = ? AND type = ?", in.UserId, in.TargetId, in.ChatType).
			Update("unread_count", 0).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		l.Logger.Errorf("标记已读失败 (事务处理失败): %v", err)
		return &chat.MarkMessageReadResponse{Success: false, Message: "标记已读失败"}, nil
	}

	// 3. 发送 Kafka 事件通知对方
	event := &mq.MessageEvent{
		Type:              mq.EventMessageRead,
		FromUserID:        in.TargetId,
		ToUserID:          in.UserId,
		ChatType:          int64(in.ChatType),
		GroupID:           in.TargetId,
		LastReadMessageID: in.LastReadMessageId,
		CreateAt:          time.Now().Unix(),
	}

	if err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, event); err != nil {
		l.Logger.Errorf("发送已读回执事件到 Kafka 失败: %v", err)
	}

	return &chat.MarkMessageReadResponse{Success: true, Message: "标记已读成功"}, nil
}
