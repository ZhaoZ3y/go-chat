package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInviteFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInviteFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInviteFriendsLogic {
	return &GetInviteFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInviteFriendsLogic) GetInviteFriends(in *group.GetInviteFriendsRequest) (*group.GetInviteFriendsResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GetInviteFriendsResponse{}, nil
}
