package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeGroupLogic {
	return &ChangeGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeGroupLogic) ChangeGroup(in *contact.ContactChangeGroupRequest) (*contact.ContactChangeGroupResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactChangeGroupResponse{}, nil
}
