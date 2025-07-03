package logic

import (
	"IM/pkg/model"
	"context"
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

	if operatorMember.Role != int64(group.MemberRole_ROLE_OWNER) {
		return &group.SetMemberRoleResponse{Success: false, Message: "只有群主才能设置管理员"}, nil
	}
	if targetMember.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.SetMemberRoleResponse{Success: false, Message: "不能更改群主角色"}, nil
	}
	if targetMember.Role == int64(in.Role) {
		return &group.SetMemberRoleResponse{Success: true, Message: "角色未发生变化"}, nil
	}

	// 新增：如果要设置角色为管理员，先检查当前管理员数量是否超过20
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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新目标成员的角色
	if err := tx.Model(&model.GroupMembers{}).Where("id = ?", targetMember.Id).Update("role", in.Role).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("SetMemberRole: update member role failed: %v", err)
		return &group.SetMemberRoleResponse{Success: false, Message: "更新角色失败"}, nil
	}

	roleName := "普通成员"
	if in.Role == group.MemberRole_ROLE_ADMIN {
		roleName = "管理员"
	}
	notificationMessage := fmt.Sprintf("您在该群的角色已被群主'%s'设置为'%s'", operatorMember.Nickname, roleName)

	notification := &model.GroupNotification{
		Type:         int64(group.NotificationType_NOTIFY_MEMBER_ROLE_CHANGED),
		GroupId:      in.GroupId,
		OperatorId:   in.OperatorId,
		TargetUserId: in.UserId,
		Message:      notificationMessage,
	}
	if err := tx.Create(notification).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("SetMemberRole: create notification failed: %v", err)
		return &group.SetMemberRoleResponse{Success: false, Message: "创建通知失败"}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("SetMemberRole: commit transaction failed: %v", err)
		return &group.SetMemberRoleResponse{Success: false, Message: "处理失败"}, nil
	}

	// TODO： 异步通知群主和管理员，类型为 NOTIFY_MEMBER_ROLE_CHANGED

	return &group.SetMemberRoleResponse{
		Success: true,
		Message: "成员角色设置成功",
	}, nil
}
