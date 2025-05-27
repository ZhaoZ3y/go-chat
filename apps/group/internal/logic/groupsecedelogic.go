package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSecedeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupSecedeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSecedeLogic {
	return &GroupSecedeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupSecedeLogic) GroupSecede(in *group.GroupSecedeRequest) (*group.GroupSecedeResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupSecedeResponse{}, nil
}
