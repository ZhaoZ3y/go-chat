package logic

import (
	"IM/pkg/model"
	"context"
	"time"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建群组
func (l *CreateGroupLogic) CreateGroup(in *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
	initialMemberIDs := make(map[int64]struct{})
	for _, userId := range in.MemberIds {
		if userId != in.OwnerId {
			initialMemberIDs[userId] = struct{}{}
		}
	}
	totalMemberCount := 1 + len(initialMemberIDs)

	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		logx.Errorf("CreateGroup: begin transaction failed, error: %v", tx.Error)
		return nil, tx.Error
	}
	// 确保在函数退出时，如果事务未提交，则回滚
	defer tx.Rollback()

	groupModel := &model.Groups{
		Name:           in.Name,
		Description:    in.Description,
		Avatar:         in.Avatar,
		OwnerId:        in.OwnerId,
		MemberCount:    int64(totalMemberCount), // 直接使用计算好的成员总数
		MaxMemberCount: 500,                     // 默认最大成员数，可根据需求调整
		Status:         int64(group.GroupStatus_GROUP_STATUS_NORMAL),
	}

	if err := tx.Create(groupModel).Error; err != nil {
		logx.Errorf("CreateGroup: create group failed, error: %v", err)
		return &group.CreateGroupResponse{Success: false, Message: "创建群组记录失败"}, nil
	}

	now := time.Now().Unix()
	membersToCreate := make([]*model.GroupMembers, 0, totalMemberCount)

	var ownerUser model.User
	if err := tx.Where("id = ?", in.OwnerId).First(&ownerUser).Error; err != nil {
		logx.Errorf("CreateGroup: get owner info failed, error: %v", err)
		return &group.CreateGroupResponse{Success: false, Message: "获取群主信息失败"}, nil
	}

	// 添加群主
	ownerMember := &model.GroupMembers{
		GroupId:  groupModel.Id,
		UserId:   in.OwnerId,
		Role:     int64(group.MemberRole_ROLE_OWNER), // 群主角色
		Nickname: ownerUser.Nickname,
		Status:   int64(group.MemberStatus_MEMBER_STATUS_NORMAL), // 修正：新成员状态应为正常
		JoinTime: now,
	}
	membersToCreate = append(membersToCreate, ownerMember)

	// 添加其他初始成员
	for userId := range initialMemberIDs {
		var user model.User
		if err := tx.Where("id = ?", userId).First(&user).Error; err != nil {
			if err.Error() == "record not found" {
				logx.Errorf("CreateGroup: user %d not found, skipping", userId)
				continue // 如果用户不存在，跳过该成员
			}
			logx.Errorf("CreateGroup: get member info failed for user %d, error: %v", userId, err)
			return &group.CreateGroupResponse{Success: false, Message: "获取群组成员信息失败"}, nil
		}
		member := &model.GroupMembers{
			GroupId:  groupModel.Id,
			UserId:   userId,
			Nickname: user.Nickname,
			Role:     int64(group.MemberRole_ROLE_MEMBER),
			Status:   int64(group.MemberStatus_MEMBER_STATUS_NORMAL),
			JoinTime: now,
		}
		membersToCreate = append(membersToCreate, member)
	}

	if len(membersToCreate) > 0 {
		if err := tx.Create(&membersToCreate).Error; err != nil {
			logx.Errorf("CreateGroup: batch create group members failed, error: %v", err)
			return &group.CreateGroupResponse{Success: false, Message: "添加群组成员失败"}, nil
		}
	}

	if err := tx.Commit().Error; err != nil {
		logx.Errorf("CreateGroup: commit transaction failed, error: %v", err)
		return &group.CreateGroupResponse{Success: false, Message: "创建群组事务提交失败"}, nil
	}

	return &group.CreateGroupResponse{
		GroupId: groupModel.Id,
		Success: true,
		Message: "创建群组成功",
	}, nil
}
