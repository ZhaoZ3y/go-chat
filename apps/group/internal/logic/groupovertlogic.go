package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupOvertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupOvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupOvertLogic {
	return &GroupOvertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupOvertLogic) GroupOvert(in *group.GroupOvertRequest) (*group.GroupOvertResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupOvertResponse{}, nil
}
