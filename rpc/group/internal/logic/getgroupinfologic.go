package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组信息
func (l *GetGroupInfoLogic) GetGroupInfo(in *group.GetGroupInfoRequest) (*group.GetGroupInfoResponse, error) {
	// 获取群组信息
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND status = 1", in.GroupId).First(&groupInfo).Error; err != nil {
		return &group.GetGroupInfoResponse{}, nil
	}

	// 获取用户在群组中的信息
	var userMember model.GroupMembers
	l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&userMember)

	return &group.GetGroupInfoResponse{
		GroupInfo: &group.Group{
			Id:             groupInfo.Id,
			Name:           groupInfo.Name,
			Description:    groupInfo.Description,
			Avatar:         groupInfo.Avatar,
			OwnerId:        groupInfo.OwnerId,
			MemberCount:    int32(groupInfo.MemberCount),
			MaxMemberCount: int32(groupInfo.MaxMemberCount),
			Status:         int32(groupInfo.Status),
			CreateAt:       groupInfo.CreateAt,
			UpdateAt:       groupInfo.UpdateAt,
		},
		UserMemberInfo: &group.GroupMember{
			Id:       userMember.Id,
			GroupId:  userMember.GroupId,
			UserId:   userMember.UserId,
			Role:     int32(userMember.Role),
			Nickname: userMember.Nickname,
			Status:   int32(userMember.Status),
			JoinTime: userMember.JoinTime,
			UpdateAt: userMember.UpdateAt,
		},
	}, nil
}
