package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInviteToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteToGroupLogic {
	return &InviteToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 邀请加入群组
func (l *InviteToGroupLogic) InviteToGroup(in *group.InviteToGroupRequest) (*group.InviteToGroupResponse, error) {
	// todo: add your logic here and delete this line

	return &group.InviteToGroupResponse{}, nil
}
