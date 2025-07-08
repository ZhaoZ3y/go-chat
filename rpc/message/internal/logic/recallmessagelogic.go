package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	_const "IM/pkg/utils/const"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecallMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecallMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallMessageLogic {
	return &RecallMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤回消息
func (l *RecallMessageLogic) RecallMessage(in *chat.RecallMessageRequest) (*chat.RecallMessageResponse, error) {
	var message model.Messages
	err := l.svcCtx.DB.Where("id = ? AND from_user_id = ?", in.MessageId, in.UserId).First(&message).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &chat.RecallMessageResponse{Success: false, Message: "消息不存在或无权限撤回"}, nil
	} else if err != nil {
		l.Logger.Errorf("查询待撤回消息失败: %v", err)
		return nil, err
	}

	if time.Now().Unix()-message.CreateAt > 120 { // 120秒 = 2分钟
		return &chat.RecallMessageResponse{Success: false, Message: "超过2分钟，无法撤回"}, nil
	}

	if message.Status == _const.MsgRecall {
		return &chat.RecallMessageResponse{Success: true, Message: "消息已撤回"}, nil
	}

	// 4. 在更新数据库前，先确定接收者列表
	var recipients []int64
	if message.ChatType == _const.ChatTypePrivate {
		// 私聊：通知的接收者是对方和自己（用于多端同步）
		recipients = []int64{message.FromUserId, message.ToUserId}
	} else {
		// 查询群内所有未被删除的成员
		if err := l.svcCtx.DB.Model(&model.GroupMembers{}).Where("group_id = ? AND deleted_at IS NULL", message.GroupId).Pluck("user_id", &recipients).Error; err != nil {
			l.Logger.Errorf("查询群 %d 成员失败: %v", message.GroupId, err)
			return nil, err
		}
	}

	// 使用 GORM 的 .Updates() 方法只更新指定字段
	err = l.svcCtx.DB.Model(&message).Updates(map[string]interface{}{
		"status":  _const.MsgRecall,
		"content": "", // 撤回后内容清空
		"extra":   "", // 额外信息也建议清空
	}).Error

	if err != nil {
		l.Logger.Errorf("更新消息为撤回状态失败: %v", err)
		return &chat.RecallMessageResponse{Success: false, Message: "撤回消息失败"}, nil
	}

	recallEvent := &mq.RichMessageEvent{
		MessageEvent: mq.MessageEvent{
			Type:       mq.EventMessageRecall,
			MessageID:  message.Id,
			FromUserID: message.FromUserId,
			ToUserID:   message.ToUserId, // 私聊时仍保留
			GroupID:    message.GroupId,  // 群聊时仍保留
			ChatType:   message.ChatType,
			CreateAt:   time.Now().Unix(),
		},
		Recipients: recipients, // 附上完整的接收者列表
	}

	go func() {
		if err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, recallEvent); err != nil {
			l.Logger.Errorf("发送消息撤回事件到 Kafka 失败: %v", err)
		}
	}()

	return &chat.RecallMessageResponse{
		Success: true,
		Message: "撤回成功",
	}, nil
}
