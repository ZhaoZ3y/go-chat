package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyListLogic {
	return &ApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyListLogic) ApplyList(in *contact.ContactApplyListRequest) (*contact.ContactApplyListResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactApplyListResponse{}, nil
}
