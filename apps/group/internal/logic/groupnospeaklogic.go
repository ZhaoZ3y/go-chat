package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupNoSpeakLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupNoSpeakLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupNoSpeakLogic {
	return &GroupNoSpeakLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupNoSpeakLogic) GroupNoSpeak(in *group.GroupNoSpeakRequest) (*group.GroupNoSpeakResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupNoSpeakResponse{}, nil
}
