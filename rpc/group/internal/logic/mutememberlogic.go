package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MuteMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMuteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteMemberLogic {
	return &MuteMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁言群组成员
func (l *MuteMemberLogic) MuteMember(in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	// todo: add your logic here and delete this line

	return &group.MuteMemberResponse{}, nil
}
