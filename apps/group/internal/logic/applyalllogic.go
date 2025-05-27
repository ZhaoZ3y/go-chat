package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAllLogic {
	return &ApplyAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyAllLogic) ApplyAll(in *group.ApplyAllRequest) (*group.ApplyAllResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyAllResponse{}, nil
}
