package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHandoverLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupHandoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHandoverLogic {
	return &GroupHandoverLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupHandoverLogic) GroupHandover(in *group.GroupHandoverRequest) (*group.GroupHandoverResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupHandoverResponse{}, nil
}
