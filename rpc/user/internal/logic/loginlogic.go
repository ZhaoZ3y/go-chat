package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/encrypt"
	"IM/pkg/utils/jwt"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

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

// 用户登录
func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名和密码不能为空")
	}

	var userModel model.User
	err := l.svcCtx.DB.Where("username = ? AND deleted_at = 0", in.Username).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, status.Error(codes.Internal, "系统错误")
	}

	if encrypt.EncryptPassword(in.Password) != userModel.Password {
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}
	if userModel.Status != 1 {
		return nil, status.Error(codes.PermissionDenied, "用户已被禁用")
	}

	now := time.Now().Unix()
	_ = l.svcCtx.DB.Model(&userModel).Updates(map[string]interface{}{
		"last_login_at": now,
		"update_at":     now,
	}).Error

	accessToken, err := jwt.GenerateAccessToken(userModel.Id, userModel.Username)
	if err != nil {
		l.Logger.Errorf("生成 access token 失败: %v", err)
		return nil, status.Error(codes.Internal, "生成 Token 失败")
	}

	refreshToken, err := jwt.GenerateRefreshToken(userModel.Id, userModel.Username)
	if err != nil {
		l.Logger.Errorf("生成 refresh token 失败: %v", err)
		return nil, status.Error(codes.Internal, "生成 Token 失败")
	}

	return &user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
