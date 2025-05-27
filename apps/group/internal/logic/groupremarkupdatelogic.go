package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupRemarkUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupRemarkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupRemarkUpdateLogic {
	return &GroupRemarkUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupRemarkUpdateLogic) GroupRemarkUpdate(in *group.GroupRemarkUpdateRequest) (*group.GroupRemarkUpdateResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupRemarkUpdateResponse{}, nil
}
