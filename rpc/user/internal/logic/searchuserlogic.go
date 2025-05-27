package logic

import (
	"context"

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

// 搜索用户
func (l *SearchUserLogic) SearchUser(in *user.SearchUserRequest) (*user.SearchUserResponse, error) {
	// todo: add your logic here and delete this line

	return &user.SearchUserResponse{}, nil
}
