package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDismissLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDismissLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDismissLogic {
	return &GroupDismissLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupDismissLogic) GroupDismiss(in *group.GroupDismissRequest) (*group.GroupDismissResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupDismissResponse{}, nil
}
