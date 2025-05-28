package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组成员列表
func (l *GetGroupMemberListLogic) GetGroupMemberList(in *group.GetGroupMemberListRequest) (*group.GetGroupMemberListResponse, error) {
	// 检查请求者是否为群成员
	var requestMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&requestMember).Error; err != nil {
		return &group.GetGroupMemberListResponse{}, nil
	}

	var members []model.GroupMembers
	var total int64

	query := l.svcCtx.DB.Where("group_id = ?", in.GroupId)
	query.Count(&total)

	offset := (in.Page - 1) * in.PageSize
	if err := query.Offset(int(offset)).Limit(int(in.PageSize)).Find(&members).Error; err != nil {
		return &group.GetGroupMemberListResponse{}, nil
	}

	var memberList []*group.GroupMember
	for _, m := range members {
		memberList = append(memberList, &group.GroupMember{
			Id:       m.Id,
			GroupId:  m.GroupId,
			UserId:   m.UserId,
			Role:     int32(m.Role),
			Nickname: m.Nickname,
			Status:   int32(m.Status),
			JoinTime: m.JoinTime,
			UpdateAt: m.UpdateAt,
		})
	}

	return &group.GetGroupMemberListResponse{
		Members: memberList,
		Total:   total,
	}, nil
}
