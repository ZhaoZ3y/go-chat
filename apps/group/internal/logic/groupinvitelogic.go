package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInviteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInviteLogic {
	return &GroupInviteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupInviteLogic) GroupInvite(in *group.GroupInviteRequest) (*group.GroupInviteResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupInviteResponse{}, nil
}
