package logic

import (
	"IM/pkg/model"
	"context"
	"strings"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserLogic) SearchUser(in *user.SearchUserReq) (*user.SearchUserResp, error) {
	var users []model.User
	var total int64

	query := l.svcCtx.DB.Model(&model.User{}).Where("status = 1")

	if in.Keyword != "" {
		keyword := "%" + strings.TrimSpace(in.Keyword) + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", keyword, keyword)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (in.Page - 1) * in.Size
	err = query.Offset(int(offset)).Limit(int(in.Size)).Find(&users).Error
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

	return &user.SearchUserResp{
		Users: result,
		Total: total,
	}, nil
}
