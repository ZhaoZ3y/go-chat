package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetMemberRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMemberRoleLogic {
	return &SetMemberRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设置群组成员角色
func (l *SetMemberRoleLogic) SetMemberRole(in *group.SetMemberRoleRequest) (*group.SetMemberRoleResponse, error) {
	// todo: add your logic here and delete this line

	return &group.SetMemberRoleResponse{}, nil
}
