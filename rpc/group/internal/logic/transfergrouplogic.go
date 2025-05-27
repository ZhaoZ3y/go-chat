package logic

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferGroupLogic {
	return &TransferGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 转让群组
func (l *TransferGroupLogic) TransferGroup(in *group.TransferGroupRequest) (*group.TransferGroupResponse, error) {
	// todo: add your logic here and delete this line

	return &group.TransferGroupResponse{}, nil
}
