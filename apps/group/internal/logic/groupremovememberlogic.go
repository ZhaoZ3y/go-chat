package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupRemoveMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupRemoveMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupRemoveMemberLogic {
	return &GroupRemoveMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupRemoveMemberLogic) GroupRemoveMember(in *group.GroupRemoveMemberRequest) (*group.GroupRemoveMemberResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupRemoveMemberResponse{}, nil
}
