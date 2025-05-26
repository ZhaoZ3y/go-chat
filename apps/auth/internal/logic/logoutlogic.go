package logic

import (
	"context"

	"IM/apps/auth/internal/svc"
	"IM/apps/auth/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 注销登录
func (l *LogoutLogic) Logout(in *auth.AuthLogoutRequest) (*auth.AuthLogoutResponse, error) {
	// todo: add your logic here and delete this line

	return &auth.AuthLogoutResponse{}, nil
}
