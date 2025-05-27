package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &chat.MarkMessageReadResponse{}, nil
}
