package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupUpdateLogic) GroupUpdate(in *contact.ContactGroupUpdateRequest) (*contact.ContactGroupUpdateResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactGroupUpdateResponse{}, nil
}
