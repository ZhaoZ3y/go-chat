package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAgreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAgreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAgreeLogic {
	return &ApplyAgreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyAgreeLogic) ApplyAgree(in *group.ApplyAgreeRequest) (*group.ApplyAgreeResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyAgreeResponse{}, nil
}
