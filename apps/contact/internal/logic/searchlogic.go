package logic

import (
	"context"

	"IM/apps/contact/internal/svc"
	"IM/apps/contact/rpc/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchLogic) Search(in *contact.ContactSearchRequest) (*contact.ContactSearchResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.ContactSearchResponse{}, nil
}
