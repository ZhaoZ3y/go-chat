package logic

import (
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendRemarkLogic {
	return &UpdateFriendRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新好友备注
func (l *UpdateFriendRemarkLogic) UpdateFriendRemark(in *friend.UpdateFriendRemarkRequest) (*friend.UpdateFriendRemarkResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.UpdateFriendRemarkResponse{}, nil
}
