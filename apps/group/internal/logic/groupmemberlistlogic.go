package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberListLogic {
	return &GroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupMemberListLogic) GroupMemberList(in *group.GroupMemberListRequest) (*group.GroupMemberListResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupMemberListResponse{}, nil
}
