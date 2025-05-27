package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupOvertListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupOvertListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupOvertListLogic {
	return &GroupOvertListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupOvertListLogic) GroupOvertList(in *group.GroupOvertListRequest) (*group.GroupOvertListResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupOvertListResponse{}, nil
}
