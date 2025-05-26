package logic

import (
	"context"

	"IM/apps/auth/internal/svc"
	"IM/apps/auth/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 刷新Token
func (l *RefreshLogic) Refresh(in *auth.AuthRefreshRequest) (*auth.AuthRefreshResponse, error) {
	// todo: add your logic here and delete this line

	return &auth.AuthRefreshResponse{}, nil
}
