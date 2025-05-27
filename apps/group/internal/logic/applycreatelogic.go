package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateLogic {
	return &ApplyCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 入群申请操作
func (l *ApplyCreateLogic) ApplyCreate(in *group.ApplyCreateRequest) (*group.ApplyCreateResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyCreateResponse{}, nil
}
