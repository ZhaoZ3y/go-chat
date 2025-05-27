package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组成员列表
func (l *GetGroupMemberListLogic) GetGroupMemberList(in *group.GetGroupMemberListRequest) (*group.GetGroupMemberListResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GetGroupMemberListResponse{}, nil
}
