package logic

import (
	"IM/pkg/model"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

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
	// 1. 权限验证：检查请求者是否为该群组成员
	var count int64
	err := l.svcCtx.DB.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Count(&count).Error
	if err != nil {
		logx.Errorf("GetGroupMemberList: check membership failed, group_id: %d, user_id: %d, error: %v", in.GroupId, in.UserId, err)
		return nil, status.Error(codes.Internal, "数据库查询失败")
	}
	if count == 0 {
		logx.Errorf("GetGroupMemberList: permission denied, user %d is not a member of group %d", in.UserId, in.GroupId)
		return nil, status.Errorf(codes.PermissionDenied, "您不是该群组成员，无权查看成员列表")
	}

	var memberModels []model.GroupMembers
	// 推荐排序：按角色降序（群主、管理员、成员），然后按加入时间升序
	err = l.svcCtx.DB.
		Where("group_id = ?", in.GroupId).
		Order("role DESC, CONVERT(nickname USING gbk) ASC").
		Find(&memberModels).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.GetGroupMemberListResponse{Members: []*group.GroupMember{}, Total: 0}, nil
		}
		logx.Errorf("GetGroupMemberList: find members failed, group_id: %d, error: %v", in.GroupId, err)
		return nil, status.Error(codes.Internal, "获取成员列表失败")
	}

	pbMembers := make([]*group.GroupMember, 0, len(memberModels))
	for _, member := range memberModels {
		pbMember := &group.GroupMember{
			Id:       member.Id,
			GroupId:  member.GroupId,
			UserId:   member.UserId,
			Role:     group.MemberRole(member.Role),
			Nickname: member.Nickname,
			Status:   group.MemberStatus(member.Status),
			JoinTime: member.JoinTime,
			UpdateAt: member.UpdateAt,
		}
		pbMembers = append(pbMembers, pbMember)
	}

	return &group.GetGroupMemberListResponse{
		Members: pbMembers,
		Total:   int64(len(pbMembers)),
	}, nil
}
