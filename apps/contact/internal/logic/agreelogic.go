package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type AgreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAgreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgreeLogic {
	return &AgreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AgreeLogic) Agree(in *contact.ContactAgreeRequest) (*contact.ContactAgreeResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactAgreeResponse{}, nil
}
