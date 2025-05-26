package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectLogic {
	return &RejectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectLogic) Reject(in *contact.ContactRejectRequest) (*contact.ContactRejectResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactRejectResponse{}, nil
}
