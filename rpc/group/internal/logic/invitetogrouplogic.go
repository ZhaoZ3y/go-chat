package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInviteToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteToGroupLogic {
	return &InviteToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 邀请加入群组
func (l *InviteToGroupLogic) InviteToGroup(in *group.InviteToGroupRequest) (*group.InviteToGroupResponse, error) {
	if in.GroupId == 0 || in.InviterId == 0 || len(in.UserIds) == 0 {
		return &group.InviteToGroupResponse{Success: false, Message: "参数错误：群组ID、邀请人ID和被邀请人列表不能为空"}, nil
	}

	// 校验群组是否存在
	var targetGroup model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", in.GroupId, group.GroupStatus_GROUP_STATUS_NORMAL).
		First(&targetGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.InviteToGroupResponse{Success: false, Message: "群组不存在或已解散"}, nil
		}
		l.Logger.Errorf("InviteToGroup: 查询群组失败, groupID: %d, error: %v", in.GroupId, err)
		return &group.InviteToGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 校验邀请人是否是群成员
	var inviterMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.InviterId).
		First(&inviterMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.InviteToGroupResponse{Success: false, Message: "您不是该群成员，无法邀请他人"}, nil
		}
		l.Logger.Errorf("InviteToGroup: 查询邀请人失败, error: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "查询邀请人信息失败"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ✅ 校验被邀请用户是否存在
	var existingUsers []model.User
	if err := tx.Model(&model.User{}).Where("id IN ?", in.UserIds).Find(&existingUsers).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("InviteToGroup: 查询用户失败: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "查询被邀请用户失败"}, nil
	}

	existingUserMap := make(map[int64]struct{}, len(existingUsers))
	for _, user := range existingUsers {
		existingUserMap[user.Id] = struct{}{}
	}

	nonExistUsers := make([]int64, 0)
	for _, uid := range in.UserIds {
		if _, ok := existingUserMap[uid]; !ok {
			nonExistUsers = append(nonExistUsers, uid)
		}
	}
	if len(nonExistUsers) > 0 {
		tx.Rollback()
		msg := fmt.Sprintf("以下用户ID不存在：%v", nonExistUsers)
		l.Logger.Errorf("InviteToGroup: 用户不存在: %v", nonExistUsers)
		return &group.InviteToGroupResponse{
			Success: false,
			Message: msg,
		}, nil
	}

	// 查找已在群中的人
	var existingMembers []model.GroupMembers
	tx.Where("group_id = ? AND user_id IN ?", in.GroupId, in.UserIds).Find(&existingMembers)
	existingMemberMap := make(map[int64]struct{}, len(existingMembers))
	for _, m := range existingMembers {
		existingMemberMap[m.UserId] = struct{}{}
	}

	// 查找已有入群申请的用户
	var pendingApplications []model.JoinGroupApplications
	tx.Where("to_group_id = ? AND from_user_id IN ? AND status = ?", in.GroupId, in.UserIds, group.ApplicationStatus_PENDING).
		Find(&pendingApplications)
	pendingApplicationMap := make(map[int64]struct{}, len(pendingApplications))
	for _, app := range pendingApplications {
		pendingApplicationMap[app.FromUserId] = struct{}{}
	}

	isAdmin := inviterMember.Role == int64(group.MemberRole_ROLE_OWNER) || inviterMember.Role == int64(group.MemberRole_ROLE_ADMIN)
	failedUserIDs := make([]int64, 0)
	var message string

	if isAdmin {
		// ✅ 管理员或群主，直接加入群组
		membersToCreate := make([]*model.GroupMembers, 0)
		for _, uid := range in.UserIds {
			if _, exists := existingMemberMap[uid]; exists {
				failedUserIDs = append(failedUserIDs, uid)
				continue
			}

			var user model.User
			if err := tx.First(&user, uid).Error; err != nil {
				failedUserIDs = append(failedUserIDs, uid)
				continue
			}

			membersToCreate = append(membersToCreate, &model.GroupMembers{
				GroupId:  in.GroupId,
				UserId:   uid,
				Nickname: user.Nickname,
				Role:     int64(group.MemberRole_ROLE_MEMBER),
				Status:   int64(group.MemberStatus_MEMBER_STATUS_NORMAL),
				JoinTime: time.Now().Unix(),
			})
		}

		if len(membersToCreate) > 0 {
			if err := tx.Create(&membersToCreate).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("InviteToGroup: 添加成员失败: %v", err)
				return &group.InviteToGroupResponse{Success: false, Message: "添加新成员失败"}, nil
			}
			if err := tx.Model(&targetGroup).UpdateColumn("member_count", gorm.Expr("member_count + ?", len(membersToCreate))).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("InviteToGroup: 更新群成员数失败: %v", err)
				return &group.InviteToGroupResponse{Success: false, Message: "更新群成员数失败"}, nil
			}
		}

		message = fmt.Sprintf("成功邀请 %d 人加入群组，%d 人失败或已在群中。", len(membersToCreate), len(failedUserIDs))

	} else {
		// ✅ 普通成员，创建入群申请
		applicationsToCreate := make([]*model.JoinGroupApplications, 0)
		for _, uid := range in.UserIds {
			if _, exists := existingMemberMap[uid]; exists || pendingApplicationMap[uid] != struct{}{} {
				failedUserIDs = append(failedUserIDs, uid)
				continue
			}
			applicationsToCreate = append(applicationsToCreate, &model.JoinGroupApplications{
				FromUserId: uid,
				ToGroupId:  in.GroupId,
				Reason:     fmt.Sprintf("由成员 %s 邀请加入", inviterMember.Nickname),
				InviterId:  in.InviterId,
				Status:     int8(group.ApplicationStatus_PENDING),
			})
		}

		if len(applicationsToCreate) > 0 {
			if err := tx.Create(&applicationsToCreate).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("InviteToGroup: 创建申请失败: %v", err)
				return &group.InviteToGroupResponse{Success: false, Message: "创建入群申请失败"}, nil
			}
		}
		message = fmt.Sprintf("已为 %d 人发送入群申请，%d 人失败或已申请/在群中。", len(applicationsToCreate), len(failedUserIDs))
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("InviteToGroup: 提交事务失败: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "处理邀请失败"}, nil
	}

	return &group.InviteToGroupResponse{
		Success:       true,
		Message:       message,
		FailedUserIds: failedUserIDs,
	}, nil
}
