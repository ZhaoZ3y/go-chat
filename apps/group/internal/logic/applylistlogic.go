package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyListLogic {
	return &ApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyListLogic) ApplyList(in *group.ApplyListRequest) (*group.ApplyListResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyListResponse{}, nil
}
