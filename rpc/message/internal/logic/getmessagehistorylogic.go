package logic

import (
	"IM/pkg/model"
	"context"

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
	var messages []model.Messages
	query := l.svcCtx.DB.Where("status = ?", 1)

	// 根据聊天类型构建查询条件
	if in.ChatType == chat.ChatType_PRIVATE {
		query = query.Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND chat_type = ?",
			in.UserId, in.TargetId, in.TargetId, in.UserId, int8(chat.ChatType_PRIVATE))
	} else {
		query = query.Where("group_id = ? AND chat_type = ?", in.TargetId, int8(chat.ChatType_GROUP))
	}

	// 分页查询
	if in.LastMessageId > 0 {
		query = query.Where("id < ?", in.LastMessageId)
	}

	if in.Limit <= 0 {
		in.Limit = 20
	}

	err := query.Order("id DESC").Limit(int(in.Limit + 1)).Find(&messages).Error
	if err != nil {
		l.Logger.Errorf("查询消息历史失败: %v", err)
		return &chat.GetMessageHistoryResponse{}, err
	}

	// 检查是否还有更多消息
	hasMore := len(messages) > int(in.Limit)
	if hasMore {
		messages = messages[:len(messages)-1]
	}

	// 转换为proto消息格式
	var chatMessages []*chat.Message
	for i := len(messages) - 1; i >= 0; i-- { // 反转顺序，最新消息在后面
		msg := messages[i]
		chatMessages = append(chatMessages, &chat.Message{
			Id:         msg.Id,
			FromUserId: msg.FromUserId,
			ToUserId:   msg.ToUserId,
			GroupId:    msg.GroupId,
			Type:       chat.MessageType(msg.Type),
			Content:    msg.Content,
			Extra:      msg.Extra,
			CreateAt:   msg.CreateAt,
			ChatType:   chat.ChatType(msg.ChatType),
		})
	}

	return &chat.GetMessageHistoryResponse{
		Messages: chatMessages,
		HasMore:  hasMore,
	}, nil
}
