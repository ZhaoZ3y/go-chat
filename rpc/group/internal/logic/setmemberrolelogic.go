package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/chat_service"
	"context"
	"encoding/json"
	"fmt"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetMemberRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMemberRoleLogic {
	return &SetMemberRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设置群组成员角色
func (l *SetMemberRoleLogic) SetMemberRole(in *group.SetMemberRoleRequest) (*group.SetMemberRoleResponse, error) {
	if in.GroupId == 0 || in.OperatorId == 0 || in.UserId == 0 {
		return &group.SetMemberRoleResponse{Success: false, Message: "参数错误"}, nil
	}
	if in.OperatorId == in.UserId {
		return &group.SetMemberRoleResponse{Success: false, Message: "不能设置自己的角色"}, nil
	}
	if in.Role != group.MemberRole_ROLE_ADMIN && in.Role != group.MemberRole_ROLE_MEMBER {
		return &group.SetMemberRoleResponse{Success: false, Message: "只能将成员设置为管理员或普通成员"}, nil
	}

	// 获取操作员和目标成员信息
	var members []model.GroupMembers
	l.svcCtx.DB.Where("group_id = ? AND user_id IN ?", in.GroupId, []int64{in.OperatorId, in.UserId}).Find(&members)

	var operatorMember, targetMember *model.GroupMembers
	for i := range members {
		if members[i].UserId == in.OperatorId {
			operatorMember = &members[i]
		}
		if members[i].UserId == in.UserId {
			targetMember = &members[i]
		}
	}

	if operatorMember == nil {
		return &group.SetMemberRoleResponse{Success: false, Message: "您不是该群成员，无权操作"}, nil
	}
	if targetMember == nil {
		return &group.SetMemberRoleResponse{Success: false, Message: "目标用户不是该群成员"}, nil
	}

	var targetGroup model.Groups
	l.svcCtx.DB.First(&targetGroup, in.GroupId)

	if operatorMember.Role != int64(group.MemberRole_ROLE_OWNER) {
		return &group.SetMemberRoleResponse{Success: false, Message: "只有群主才能设置管理员"}, nil
	}
	if targetMember.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.SetMemberRoleResponse{Success: false, Message: "不能更改群主角色"}, nil
	}
	if targetMember.Role == int64(in.Role) {
		return &group.SetMemberRoleResponse{Success: true, Message: "角色未发生变化"}, nil
	}

	if in.Role == group.MemberRole_ROLE_ADMIN {
		var adminCount int64
		err := l.svcCtx.DB.Model(&model.GroupMembers{}).
			Where("group_id = ? AND role = ?", in.GroupId, int64(group.MemberRole_ROLE_ADMIN)).
			Count(&adminCount).Error
		if err != nil {
			l.Logger.Errorf("SetMemberRole: count admin members failed: %v", err)
			return &group.SetMemberRoleResponse{Success: false, Message: "获取管理员数量失败"}, nil
		}
		if adminCount >= 20 {
			return &group.SetMemberRoleResponse{Success: false, Message: "管理员数量已达上限（20人）"}, nil
		}
	}

	tx := l.svcCtx.DB.Begin()
	defer tx.Rollback()

	if err := tx.Model(&model.GroupMembers{}).Where("id = ?", targetMember.Id).Update("role", in.Role).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("SetMemberRole: update member role failed: %v", err)
		return &group.SetMemberRoleResponse{Success: false, Message: "更新角色失败"}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("SetMemberRole: commit transaction failed: %v", err)
		return &group.SetMemberRoleResponse{Success: false, Message: "处理失败"}, nil
	}

	go l.notifyTargetUser(targetMember, operatorMember, &targetGroup, in.Role)
	go l.notifyAdminsOfRoleChange(targetMember, operatorMember, &targetGroup, in.Role)

	return &group.SetMemberRoleResponse{
		Success: true,
		Message: "成员角色设置成功",
	}, nil
}

// notifyTargetUser 异步通知被操作的用户角色已变更
func (l *SetMemberRoleLogic) notifyTargetUser(target *model.GroupMembers, operator *model.GroupMembers, groupInfo *model.Groups, newRole group.MemberRole) {
	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	roleName := "普通成员"
	if newRole == group.MemberRole_ROLE_ADMIN {
		roleName = "管理员"
	}

	title := "群内角色变更通知"
	content := fmt.Sprintf("在群聊 '%s' 中，您的角色已被群主 '%s' 设置为 '%s'。", groupInfo.Name, operator.Nickname, roleName)

	extraData := map[string]interface{}{
		"notification_type": group.NotificationType_NOTIFY_MEMBER_ROLE_CHANGED.String(),
		"group_id":          groupInfo.Id,
		"group_name":        groupInfo.Name,
		"operator_id":       operator.UserId,
		"operator_nickname": operator.Nickname,
		"new_role":          newRole,
	}
	extraJSON, _ := json.Marshal(extraData)

	err := notificationSvc.SendSystemNotification(target.UserId, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("SetMemberRole-Notify: 发送角色变更通知给目标用户失败, UserID: %d, GroupID: %d, Error: %v", target.UserId, groupInfo.Id, err)
	}
}

// notifyAdminsOfRoleChange 异步通知群主和其他管理员角色变更事件
func (l *SetMemberRoleLogic) notifyAdminsOfRoleChange(target *model.GroupMembers, operator *model.GroupMembers, groupInfo *model.Groups, newRole group.MemberRole) {
	var adminIDs []int64
	err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND role = ?", groupInfo.Id, int64(group.MemberRole_ROLE_ADMIN)).
		Pluck("user_id", &adminIDs).Error
	if err != nil {
		l.Logger.Errorf("SetMemberRole-AdminNotify: 查找管理员列表失败, groupID: %d, err: %v", groupInfo.Id, err)
		return
	}

	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	roleName := "普通成员"
	if newRole == group.MemberRole_ROLE_ADMIN {
		roleName = "管理员"
	}

	title := "群成员角色变更"
	content := fmt.Sprintf("在群聊 '%s' 中，群主 '%s' 将成员 '%s' 的角色设置为了 '%s'。", groupInfo.Name, operator.Nickname, target.Nickname, roleName)

	extraData := map[string]interface{}{
		"notification_type": group.NotificationType_NOTIFY_MEMBER_ROLE_CHANGED.String(),
		"group_id":          groupInfo.Id,
		"group_name":        groupInfo.Name,
		"operator_id":       operator.UserId,
		"operator_nickname": operator.Nickname,
		"target_user_id":    target.UserId,
		"target_nickname":   target.Nickname,
		"new_role":          newRole,
	}
	extraJSON, _ := json.Marshal(extraData)

	// 3. 批量发送
	err = notificationSvc.SendBatchNotification(adminIDs, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("SetMemberRole-AdminNotify: 批量发送管理员通知失败: %v", err)
	} else {
		l.Logger.Infof("SetMemberRole-AdminNotify: 已成功向群 %d 的 %d 位管理员发送角色变更通知。", groupInfo.Id, len(adminIDs))
	}
}
