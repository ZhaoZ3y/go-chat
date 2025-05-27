package logic

import (
	"IM/pkg/model"
	"context"
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"time"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 检查用户名是否已存在
	var existUser model.User
	err := l.svcCtx.DB.Where("username = ?", in.Username).First(&existUser).Error
	if err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 检查邮箱是否已存在
	if in.Email != "" {
		err = l.svcCtx.DB.Where("email = ?", in.Email).First(&existUser).Error
		if err == nil {
			return nil, fmt.Errorf("邮箱已被使用")
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	// 密码加密
	hasher := md5.New()
	hasher.Write([]byte(in.Password))
	hashedPassword := fmt.Sprintf("%x", hasher.Sum(nil))

	// 创建用户
	newUser := &model.User{
		Username: in.Username,
		Password: hashedPassword,
		Nickname: in.Nickname,
		Email:    in.Email,
		Status:   1,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}

	err = l.svcCtx.DB.Create(newUser).Error
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Id:       newUser.Id,
		Username: newUser.Username,
		Nickname: newUser.Nickname,
		Avatar:   newUser.Avatar,
		Email:    newUser.Email,
		CreateAt: newUser.CreateAt,
	}, nil
}
