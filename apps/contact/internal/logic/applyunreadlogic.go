package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyUnreadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyUnreadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyUnreadLogic {
	return &ApplyUnreadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyUnreadLogic) ApplyUnread(in *contact.ContactApplyUnreadRequest) (*contact.ContactApplyUnreadResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactApplyUnreadResponse{}, nil
}
