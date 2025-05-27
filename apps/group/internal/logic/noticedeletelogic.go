package logic

import (
	"context"

	"IM/apps/group/internal/svc"
	"IM/apps/group/rpc/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeDeleteLogic {
	return &NoticeDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群公告操作
func (l *NoticeDeleteLogic) NoticeDelete(in *group.NoticeDeleteRequest) (*group.NoticeDeleteResponse, error) {
	// todo: add your logic here and delete this line

	return &group.NoticeDeleteResponse{}, nil
}
