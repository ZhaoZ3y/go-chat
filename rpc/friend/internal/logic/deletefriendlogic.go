package logic

import (
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除好友
func (l *DeleteFriendLogic) DeleteFriend(in *friend.DeleteFriendRequest) (*friend.DeleteFriendResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.DeleteFriendResponse{}, nil
}
