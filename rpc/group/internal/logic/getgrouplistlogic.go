package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupListLogic {
	return &GetGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组列表
func (l *GetGroupListLogic) GetGroupList(in *group.GetGroupListRequest) (*group.GetGroupListResponse, error) {
	var groups []model.Groups
	var total int64

	// 通过群成员表查询用户所在的群组
	query := l.svcCtx.DB.Table("groups g").
		Joins("JOIN group_members gm ON g.id = gm.group_id").
		Where("gm.user_id = ? AND g.status = 1", in.UserId)

	// 计算总数
	query.Count(&total)

	// 分页查询
	if err := query.Find(&groups).Error; err != nil {
		return &group.GetGroupListResponse{}, nil
	}

	var groupList []*group.Group
	for _, g := range groups {
		groupList = append(groupList, &group.Group{
			Id:             g.Id,
			Name:           g.Name,
			Description:    g.Description,
			Avatar:         g.Avatar,
			OwnerId:        g.OwnerId,
			MemberCount:    int32(g.MemberCount),
			MaxMemberCount: int32(g.MaxMemberCount),
			Status:         int32(g.Status),
			CreateAt:       g.CreateAt,
			UpdateAt:       g.UpdateAt,
		})
	}

	return &group.GetGroupListResponse{
		Groups: groupList,
		Total:  total,
	}, nil
}
