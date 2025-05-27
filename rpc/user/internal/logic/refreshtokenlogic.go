package logic

import (
	"IM/pkg/utils/jwt"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 刷新Token
func (l *RefreshTokenLogic) RefreshToken(in *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error) {
	claims, err := jwt.ParseRefreshToken(in.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "无效的 refresh token")
	}

	accessToken, err := jwt.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "生成 access token 失败")
	}

	refreshToken, err := jwt.GenerateRefreshToken(claims.UserID, claims.Username)
	if err != nil {
		l.Logger.Errorf("生成 refresh token 失败: %v", err)
		return nil, status.Error(codes.Internal, "生成 Token 失败")
	}

	return &user.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
