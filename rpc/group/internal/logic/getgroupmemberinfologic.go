package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberInfoLogic {
	return &GetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组成员信息
func (l *GetGroupMemberInfoLogic) GetGroupMemberInfo(in *group.GetGroupMemberInfoRequest) (*group.GetGroupMemberInfoResponse, error) {
	var user model.User
	err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&user).Error
	if err != nil {
		logx.Errorf("GetGroupMemberInfo: find user failed, user_id: %d, error: %v", in.UserId, err)
		return nil, err
	}

	var member model.GroupMembers
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&member).Error
	if err != nil {
		logx.Errorf("GetGroupMemberInfo: find group member failed, group_id: %d, user_id: %d, error: %v", in.GroupId, in.UserId, err)
		return nil, err
	}
	groupMemberInfo := &group.GroupMemberInfo{
		MemberInfo: &group.GroupMember{
			UserId:   member.UserId,
			GroupId:  member.GroupId,
			Role:     group.MemberRole(member.Role),
			Nickname: member.Nickname,
			Status:   group.MemberStatus(member.Status),
			JoinTime: member.JoinTime,
		},
		UserInfo: &group.User{
			Id:       user.Id,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
		},
	}

	return &group.GetGroupMemberInfoResponse{
		Info: groupMemberInfo,
	}, nil
}
