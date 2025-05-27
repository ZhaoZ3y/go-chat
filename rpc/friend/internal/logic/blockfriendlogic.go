package logic

import (
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockFriendLogic {
	return &BlockFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 拉黑好友
func (l *BlockFriendLogic) BlockFriend(in *friend.BlockFriendRequest) (*friend.BlockFriendResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.BlockFriendResponse{}, nil
}
