package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeEditLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeEditLogic {
	return &NoticeEditLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NoticeEditLogic) NoticeEdit(in *group.NoticeEditRequest) (*group.NoticeEditResponse, error) {
	// todo: add your logic here and delete this line

	return &group.NoticeEditResponse{}, nil
}
