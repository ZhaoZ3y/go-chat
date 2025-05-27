package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupInfoLogic {
	return &UpdateGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新群组信息
func (l *UpdateGroupInfoLogic) UpdateGroupInfo(in *group.UpdateGroupInfoRequest) (*group.UpdateGroupInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &group.UpdateGroupInfoResponse{}, nil
}
