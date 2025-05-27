package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByIdsLogic {
	return &GetUsersByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersByIdsLogic) GetUsersByIds(in *user.GetUsersByIdsReq) (*user.GetUsersByIdsResp, error) {
	if len(in.UserIds) == 0 {
		return &user.GetUsersByIdsResp{Users: []*user.GetUserInfoResp{}}, nil
	}

	var users []model.User
	err := l.svcCtx.DB.Where("id IN ?", in.UserIds).Find(&users).Error
	if err != nil {
		return nil, err
	}

	var result []*user.GetUserInfoResp
	for _, u := range users {
		result = append(result, &user.GetUserInfoResp{
			Id:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Email:    u.Email,
			Status:   u.Status,
			CreateAt: u.CreateAt,
			UpdateAt: u.UpdateAt,
		})
	}

	return &user.GetUsersByIdsResp{Users: result}, nil
}
