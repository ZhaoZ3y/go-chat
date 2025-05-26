package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteLogic {
	return &GroupDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupDeleteLogic) GroupDelete(in *contact.ContactGroupDeleteRequest) (*contact.ContactGroupDeleteResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactGroupDeleteResponse{}, nil
}
