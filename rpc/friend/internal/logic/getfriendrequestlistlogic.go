package logic

import (
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendRequestListLogic {
	return &GetFriendRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取好友申请列表
func (l *GetFriendRequestListLogic) GetFriendRequestList(in *friend.GetFriendRequestListRequest) (*friend.GetFriendRequestListResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.GetFriendRequestListResponse{}, nil
}
