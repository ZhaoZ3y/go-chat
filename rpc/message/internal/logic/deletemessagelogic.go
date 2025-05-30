package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除消息
func (l *DeleteMessageLogic) DeleteMessage(in *chat.DeleteMessageRequest) (*chat.DeleteMessageResponse, error) {
	// 检查消息是否存在且属于当前用户
	var message model.Messages
	err := l.svcCtx.DB.Where("id = ? AND from_user_id = ?", in.MessageId, in.UserId).First(&message).Error
	if err != nil {
		l.Logger.Errorf("消息不存在或无权限删除: %v", err)
		return &chat.DeleteMessageResponse{
			Success: false,
			Message: "消息不存在或无权限删除",
		}, nil
	}

	// 软删除消息
	err = l.svcCtx.DB.Delete(&message).Error
	if err != nil {
		l.Logger.Errorf("删除消息失败: %v", err)
		return &chat.DeleteMessageResponse{
			Success: false,
			Message: "删除消息失败",
		}, nil
	}

	return &chat.DeleteMessageResponse{
		Success: true,
		Message: "删除成功",
	}, nil
}
