package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeleteLogic {
	return &ApplyDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyDeleteLogic) ApplyDelete(in *group.ApplyDeleteRequest) (*group.ApplyDeleteResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyDeleteResponse{}, nil
}
