package logic

import (
	"IM/pkg/model"
	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"
	"context"
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"time"

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

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 密码加密
	hasher := md5.New()
	hasher.Write([]byte(in.Password))
	hashedPassword := fmt.Sprintf("%x", hasher.Sum(nil))

	// 查找用户
	var userInfo model.User
	err := l.svcCtx.DB.Where("username = ? AND password = ?", in.Username, hashedPassword).First(&userInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		return nil, err
	}

	// 检查用户状态
	if userInfo.Status != 1 {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 更新最后登录时间
	now := time.Now()
	l.svcCtx.DB.Model(&userInfo).Update("last_login_at", now)

	return &user.LoginResp{
		Id:       userInfo.Id,
		Username: userInfo.Username,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.Avatar,
		Email:    userInfo.Email,
		Status:   userInfo.Status,
		CreateAt: userInfo.CreateAt,
	}, nil
}
