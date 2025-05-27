package logic

import (
	"context"

	"IM/rpc/message/chat"
	"IM/rpc/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationListLogic {
	return &GetConversationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话列表
func (l *GetConversationListLogic) GetConversationList(in *chat.GetConversationListRequest) (*chat.GetConversationListResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.GetConversationListResponse{}, nil
}
