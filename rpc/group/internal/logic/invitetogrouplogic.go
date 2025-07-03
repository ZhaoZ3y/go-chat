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
	// 1. 前置参数校验
	if in.GroupId == 0 || in.InviterId == 0 || len(in.UserIds) == 0 {
		return &group.InviteToGroupResponse{Success: false, Message: "参数错误：群组ID、邀请人ID和被邀请人列表不能为空"}, nil
	}

	var targetGroup model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", in.GroupId, group.GroupStatus_GROUP_STATUS_NORMAL).First(&targetGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.InviteToGroupResponse{Success: false, Message: "群组不存在或已解散"}, nil
		}
		l.Logger.Errorf("InviteToGroup: find group failed, GroupID: %d, Error: %v", in.GroupId, err)
		return &group.InviteToGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	var inviterMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.InviterId).First(&inviterMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.InviteToGroupResponse{Success: false, Message: "您不是该群成员，无法邀请他人"}, nil
		}
		l.Logger.Errorf("InviteToGroup: find inviter member info failed, Error: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "查询邀请人信息失败"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingMembers []model.GroupMembers
	tx.Where("group_id = ? AND user_id IN ?", in.GroupId, in.UserIds).Find(&existingMembers)
	existingMemberMap := make(map[int64]struct{})
	for _, member := range existingMembers {
		existingMemberMap[member.UserId] = struct{}{}
	}

	var pendingApplications []model.JoinGroupApplications
	tx.Where("to_group_id = ? AND from_user_id IN ? AND status = ?", in.GroupId, in.UserIds, group.ApplicationStatus_PENDING).Find(&pendingApplications)
	pendingApplicationMap := make(map[int64]struct{})
	for _, app := range pendingApplications {
		pendingApplicationMap[app.FromUserId] = struct{}{}
	}

	failedUserIDs := make([]int64, 0)
	var message string

	isInviterAdmin := inviterMember.Role == int64(group.MemberRole_ROLE_OWNER) || inviterMember.Role == int64(group.MemberRole_ROLE_ADMIN)

	if isInviterAdmin {
		// --- 场景一：管理员或群主邀请，直接入群 ---
		membersToCreate := make([]*model.GroupMembers, 0)
		for _, inviteeID := range in.UserIds {
			if _, ok := existingMemberMap[inviteeID]; ok {
				failedUserIDs = append(failedUserIDs, inviteeID)
				continue
			}
			var user model.User
			if err := tx.First(&user, inviteeID).Error; err != nil {
				failedUserIDs = append(failedUserIDs, inviteeID)
				continue
			}

			membersToCreate = append(membersToCreate, &model.GroupMembers{
				GroupId:  in.GroupId,
				UserId:   inviteeID,
				Nickname: user.Nickname,
				Role:     int64(group.MemberRole_ROLE_MEMBER),
				Status:   int64(group.MemberStatus_MEMBER_STATUS_NORMAL),
				JoinTime: time.Now().Unix(),
			})
		}

		if len(membersToCreate) > 0 {
			// 批量创建成员
			if err := tx.Create(&membersToCreate).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("InviteToGroup (Admin): batch create members failed: %v", err)
				return &group.InviteToGroupResponse{Success: false, Message: "添加新成员失败"}, nil
			}
			// 批量更新群成员数
			if err := tx.Model(&targetGroup).UpdateColumn("member_count", gorm.Expr("member_count + ?", len(membersToCreate))).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("InviteToGroup (Admin): update member count failed: %v", err)
				return &group.InviteToGroupResponse{Success: false, Message: "更新群成员数失败"}, nil
			}
		}
		message = fmt.Sprintf("成功邀请 %d 人加入群组，%d 人失败或已在群组中。", len(membersToCreate), len(failedUserIDs))

		// TODO：后续添加消息队列实现异步通知其他群主和管理员，该高级成员邀请用户进入群聊，类型为 NOTIFY_MEMBER_BE_INVITED
		//
	} else {
		// --- 场景二：普通成员邀请，生成入群申请 ---
		applicationsToCreate := make([]*model.JoinGroupApplications, 0)
		for _, inviteeID := range in.UserIds {
			if _, ok := existingMemberMap[inviteeID]; ok { // 已是成员
				failedUserIDs = append(failedUserIDs, inviteeID)
				continue
			}
			if _, ok := pendingApplicationMap[inviteeID]; ok { // 已有待处理申请
				failedUserIDs = append(failedUserIDs, inviteeID)
				continue
			}

			applicationsToCreate = append(applicationsToCreate, &model.JoinGroupApplications{
				FromUserId: inviteeID,
				ToGroupId:  in.GroupId,
				Reason:     fmt.Sprintf("由成员 %s 邀请加入", inviterMember.Nickname),
				InviterId:  in.InviterId,
				Status:     int8(group.ApplicationStatus_PENDING),
			})
		}

		if len(applicationsToCreate) > 0 {
			if err := tx.Create(&applicationsToCreate).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("InviteToGroup (Member): batch create applications failed: %v", err)
				return &group.InviteToGroupResponse{Success: false, Message: "创建入群申请失败"}, nil
			}
		}
		message = fmt.Sprintf("已为 %d 人发送入群申请，等待管理员审核。%d 人失败或已有申请。", len(applicationsToCreate), len(failedUserIDs))

		//TODO：发送通知给管理员和群主，类型为 NOTIFY_MEMBER_APPLY_JOIN
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("InviteToGroup: commit transaction failed: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "处理邀请失败"}, nil
	}

	return &group.InviteToGroupResponse{
		Success:       true,
		Message:       message,
		FailedUserIds: failedUserIDs,
	}, nil
}
