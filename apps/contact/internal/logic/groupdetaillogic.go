package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDetailLogic {
	return &GroupDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupDetailLogic) GroupDetail(in *contact.ContactGroupDetailRequest) (*contact.ContactGroupDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactGroupDetailResponse{}, nil
}
