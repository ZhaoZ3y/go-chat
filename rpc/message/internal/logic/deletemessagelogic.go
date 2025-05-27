package logic

import (
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
	// todo: add your logic here and delete this line

	return &chat.DeleteMessageResponse{}, nil
}
