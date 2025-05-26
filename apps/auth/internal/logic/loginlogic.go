package logic

import (
	"context"

	"IM/apps/auth/internal/svc"
	"IM/apps/auth/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登录
func (l *LoginLogic) Login(in *auth.AuthLoginRequest) (*auth.AuthLoginResponse, error) {
	// todo: add your logic here and delete this line

	return &auth.AuthLoginResponse{}, nil
}
