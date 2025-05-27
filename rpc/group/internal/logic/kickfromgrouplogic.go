package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickFromGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickFromGroupLogic {
	return &KickFromGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 踢出群组
func (l *KickFromGroupLogic) KickFromGroup(in *group.KickFromGroupRequest) (*group.KickFromGroupResponse, error) {
	// todo: add your logic here and delete this line

	return &group.KickFromGroupResponse{}, nil
}
