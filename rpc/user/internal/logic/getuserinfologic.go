package logic

import (
	"IM/pkg/model"
	"context"
	"fmt"
	"gorm.io/gorm"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	var userInfo model.User
	err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&userInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	return &user.GetUserInfoResp{
		Id:       userInfo.Id,
		Username: userInfo.Username,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.Avatar,
		Email:    userInfo.Email,
		Status:   userInfo.Status,
		CreateAt: userInfo.CreateAt,
		UpdateAt: userInfo.UpdateAt,
	}, nil
}
