package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组信息
func (l *GetGroupInfoLogic) GetGroupInfo(in *group.GetGroupInfoRequest) (*group.GetGroupInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GetGroupInfoResponse{}, nil
}
