package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDetailLogic {
	return &GroupDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupDetailLogic) GroupDetail(in *group.GroupDetailRequest) (*group.GroupDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupDetailResponse{}, nil
}
