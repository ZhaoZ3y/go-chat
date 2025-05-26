package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyLogic) Apply(in *contact.ContactApplyRequest) (*contact.ContactApplyResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactApplyResponse{}, nil
}
