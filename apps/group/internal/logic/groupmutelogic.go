package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMuteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupMuteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMuteLogic {
	return &GroupMuteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupMuteLogic) GroupMute(in *group.GroupMuteRequest) (*group.GroupMuteResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupMuteResponse{}, nil
}
