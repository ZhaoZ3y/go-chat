package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyUnreadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyUnreadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyUnreadLogic {
	return &ApplyUnreadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyUnreadLogic) ApplyUnread(in *group.ApplyUnreadRequest) (*group.ApplyUnreadResponse, error) {
	// todo: add your logic here and delete this line

	return &group.ApplyUnreadResponse{}, nil
}
