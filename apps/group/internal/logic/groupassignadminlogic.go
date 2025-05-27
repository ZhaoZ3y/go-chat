package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupAssignAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupAssignAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupAssignAdminLogic {
	return &GroupAssignAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupAssignAdminLogic) GroupAssignAdmin(in *group.GroupAssignAdminRequest) (*group.GroupAssignAdminResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupAssignAdminResponse{}, nil
}
