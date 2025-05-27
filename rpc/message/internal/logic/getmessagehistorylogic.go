package logic

import (
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
	// todo: add your logic here and delete this line

	return &chat.GetMessageHistoryResponse{}, nil
}
