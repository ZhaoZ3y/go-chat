package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSettingLogic {
	return &GroupSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupSettingLogic) GroupSetting(in *group.GroupSettingRequest) (*group.GroupSettingResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupSettingResponse{}, nil
}
