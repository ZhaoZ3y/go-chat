package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeListLogic {
	return &NoticeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NoticeListLogic) NoticeList(in *group.NoticeListRequest) (*group.NoticeListResponse, error) {
	// todo: add your logic here and delete this line

	return &group.NoticeListResponse{}, nil
}
