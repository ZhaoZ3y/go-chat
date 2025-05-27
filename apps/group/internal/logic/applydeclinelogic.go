package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyDeclineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyDeclineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeclineLogic {
	return &ApplyDeclineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyDeclineLogic) ApplyDecline(in *group.ApplyDeclineRequest) (*group.ApplyDeclineResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyDeclineResponse{}, nil
}
