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
	var groupModels []model.Groups
	err := l.svcCtx.DB.
		Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", in.UserId).
		Find(&groupModels).Error

	if err != nil {
		logx.Errorf("GetGroupList: query user's groups failed, user_id: %d, error: %v", in.UserId, err)
		return nil, err
	}

	pbGroups := make([]*group.Group, 0, len(groupModels))
	for _, gm := range groupModels {
		pbGroup := &group.Group{
			Id:             gm.Id,
			Name:           gm.Name,
			Description:    gm.Description,
			Avatar:         gm.Avatar,
			OwnerId:        gm.OwnerId,
			MemberCount:    int32(gm.MemberCount),
			MaxMemberCount: int32(gm.MaxMemberCount),
			Status:         group.GroupStatus(gm.Status),
			CreateAt:       gm.CreateAt,
			UpdateAt:       gm.UpdateAt,
		}
		pbGroups = append(pbGroups, pbGroup)
	}

	return &group.GetGroupListResponse{
		Groups: pbGroups,
		Total:  int64(len(pbGroups)),
	}, nil
}
