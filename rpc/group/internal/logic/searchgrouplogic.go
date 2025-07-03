package logic

import (
	"IM/pkg/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGroupLogic {
	return &SearchGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索群组
func (l *SearchGroupLogic) SearchGroup(in *group.SearchGroupRequest) (*group.SearchGroupResponse, error) {
	if in.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "搜索关键词不能为空")
	}

	// 构建查询
	query := l.svcCtx.DB.Model(&model.Groups{}).
		Where("deleted_at IS NULL AND status = 1").
		Where("name LIKE ? OR description LIKE ?", "%"+in.Keyword+"%", "%"+in.Keyword+"%")

	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		l.Logger.Errorf("统计群组数量失败: %v", err)
		return nil, status.Error(codes.Internal, "搜索失败")
	}

	// 获取群组列表
	var groups []model.Groups
	err = query.Select("id, name, description, avatar, owner_id, member_count, max_member_count, status, create_at, update_at").
		Order("create_at DESC").
		Find(&groups).Error
	if err != nil {
		l.Logger.Errorf("查询群组列表失败: %v", err)
		return nil, status.Error(codes.Internal, "搜索失败")
	}

	groupList := make([]*group.Group, 0, len(groups))
	for _, g := range groups {
		groupList = append(groupList, &group.Group{
			Id:             g.Id,
			Name:           g.Name,
			Description:    g.Description,
			Avatar:         g.Avatar,
			OwnerId:        g.OwnerId,
			MemberCount:    int32(g.MemberCount),
			MaxMemberCount: int32(g.MaxMemberCount),
			Status:         group.GroupStatus(g.Status),
			CreateAt:       g.CreateAt,
			UpdateAt:       g.UpdateAt,
		})
	}

	return &group.SearchGroupResponse{
		Groups: groupList,
		Total:  total,
	}, nil
}
