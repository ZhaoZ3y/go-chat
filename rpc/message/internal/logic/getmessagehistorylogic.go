package logic

import (
	"IM/pkg/model"
	_const "IM/pkg/utils/const"
	"context"
	"time"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageHistoryLogic {
	return &GetMessageHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取消息历史
func (l *GetMessageHistoryLogic) GetMessageHistory(in *chat.GetMessageHistoryRequest) (*chat.GetMessageHistoryResponse, error) {
	l.Logger.Infof("[GetMessageHistory] 查询参数: userId=%d targetId=%d chatType=%d lastMsgId=%d date=%d limit=%d",
		in.UserId, in.TargetId, in.ChatType, in.LastMessageId, in.Date, in.Limit)

	query := l.svcCtx.DB.Model(&model.Messages{}).Where("messages.status = ?", _const.MsgNormal)

	// 联表 message_user_states（可选，调试时可注释掉）
	query = query.Joins("LEFT JOIN message_user_states ON messages.id = message_user_states.message_id AND message_user_states.user_id = ?", in.UserId).
		Where("message_user_states.is_deleted IS NULL OR message_user_states.is_deleted = false")

	// 私聊或群聊过滤
	if in.ChatType == chat.ChatType_PRIVATE {
		query = query.Where("((messages.from_user_id = ? AND messages.to_user_id = ?) OR (messages.from_user_id = ? AND messages.to_user_id = ?)) AND messages.chat_type = ?",
			in.UserId, in.TargetId, in.TargetId, in.UserId, int8(chat.ChatType_PRIVATE))
	} else {
		query = query.Where("messages.group_id = ? AND messages.chat_type = ?", in.TargetId, int8(chat.ChatType_GROUP))
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}

	// 日期过滤 or 滚动加载
	if in.Date > 0 {
		t := time.Unix(in.Date, 0)
		startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
		endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location()).Unix()
		l.Logger.Infof("[GetMessageHistory] 时间过滤: %d ~ %d", startOfDay, endOfDay)
		query = query.Where("messages.create_at BETWEEN ? AND ?", startOfDay, endOfDay)
		query = query.Order("messages.id ASC")
	} else {
		if in.LastMessageId > 0 {
			query = query.Where("messages.id < ?", in.LastMessageId)
		}
		query = query.Order("messages.id DESC")
	}

	var resultMessages []model.Messages
	err := query.Limit(int(limit + 1)).Find(&resultMessages).Error
	if err != nil {
		l.Logger.Errorf("[GetMessageHistory] 查询失败: %v", err)
		return nil, err
	}

	l.Logger.Infof("[GetMessageHistory] 查询到消息数: %d", len(resultMessages))

	// 判断是否还有更多
	hasMore := len(resultMessages) > int(limit)
	if hasMore {
		resultMessages = resultMessages[:len(resultMessages)-1]
	}

	if len(resultMessages) == 0 {
		l.Logger.Info("[GetMessageHistory] 没有历史消息可返回")
		return &chat.GetMessageHistoryResponse{Messages: []*chat.Message{}, HasMore: false}, nil
	}

	// 转换格式
	chatMessages := make([]*chat.Message, 0, len(resultMessages))
	toChatMessage := func(msg model.Messages) *chat.Message {
		return &chat.Message{
			Id:         msg.Id,
			FromUserId: msg.FromUserId,
			ToUserId:   msg.ToUserId,
			GroupId:    msg.GroupId,
			Type:       chat.MessageType(msg.Type),
			Content:    msg.Content,
			Extra:      msg.Extra,
			CreateAt:   msg.CreateAt,
			ChatType:   chat.ChatType(msg.ChatType),
		}
	}

	if in.Date > 0 {
		for _, msg := range resultMessages {
			chatMessages = append(chatMessages, toChatMessage(msg))
		}
	} else {
		// 逆序返回（滚动加载）
		for i := len(resultMessages) - 1; i >= 0; i-- {
			chatMessages = append(chatMessages, toChatMessage(resultMessages[i]))
		}
	}

	l.Logger.Infof("[GetMessageHistory] 返回消息: firstId=%d createAt=%d hasMore=%v",
		chatMessages[0].Id, chatMessages[0].CreateAt, hasMore)

	return &chat.GetMessageHistoryResponse{
		Messages: chatMessages,
		HasMore:  hasMore,
	}, nil
}
