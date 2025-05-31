package logic

import (
	"IM/pkg/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	if in.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "搜索关键词不能为空")
	}

	// 设置默认分页参数
	page := in.Page
	size := in.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100 // 限制每页最大数量
	}

	offset := (page - 1) * size

	// 构建查询条件
	query := l.svcCtx.DB.Model(&model.User{}).Where("deleted_at IS NULL AND status = 1")
	query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
		"%"+in.Keyword+"%", "%"+in.Keyword+"%", "%"+in.Keyword+"%")

	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		l.Logger.Errorf("统计用户数量失败: %v", err)
		return nil, status.Error(codes.Internal, "搜索失败")
	}

	// 获取用户列表
	var users []model.User
	err = query.Select("id, username, nickname, avatar, email, status, create_at, update_at").
		Offset(int(offset)).Limit(int(size)).
		Order("create_at DESC").
		Find(&users).Error
	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, status.Error(codes.Internal, "搜索失败")
	}

	// 转换结果
	userList := make([]*user.User, 0, len(users))
	for _, u := range users {
		userList = append(userList, &user.User{
			Id:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Email:    u.Email,
			Status:   int32(u.Status),
			CreateAt: u.CreateAt,
			UpdateAt: u.UpdateAt,
		})
	}

	return &user.SearchUserResponse{
		Users: userList,
		Total: total,
	}, nil
}
