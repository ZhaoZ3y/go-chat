package logic

import (
	"IM/pkg/model"
	"context"
	"gorm.io/gorm/clause"

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
	state := model.MessageUserStates{
		MessageId: in.MessageId,
		UserId:    in.UserId,
		IsDeleted: true,
	}

	err := l.svcCtx.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_deleted"}),
	}).Create(&state).Error

	if err != nil {
		l.Logger.Errorf("为用户 %d 设置消息 %d 为已删除状态失败: %v", in.UserId, in.MessageId, err)
		return &chat.DeleteMessageResponse{
			Success: false,
			Message: "操作失败",
		}, err
	}

	l.Logger.Infof("用户 %d 成功将消息 %d 从其视图中删除", in.UserId, in.MessageId)

	return &chat.DeleteMessageResponse{
		Success: true,
		Message: "删除成功",
	}, nil
}
