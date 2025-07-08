package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/chat_service"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickFromGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickFromGroupLogic {
	return &KickFromGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 踢出群组
func (l *KickFromGroupLogic) KickFromGroup(in *group.KickFromGroupRequest) (*group.KickFromGroupResponse, error) {
	if in.GroupId == 0 || in.OperatorId == 0 || in.UserId == 0 {
		return &group.KickFromGroupResponse{Success: false, Message: "参数错误"}, nil
	}

	if in.OperatorId == in.UserId {
		return &group.KickFromGroupResponse{Success: false, Message: "不能将自己踢出群组"}, nil
	}

	// 获取群组信息
	var targetGroup model.Groups
	if err := l.svcCtx.DB.First(&targetGroup, in.GroupId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.KickFromGroupResponse{Success: false, Message: "群组不存在"}, nil
		}
		l.Logger.Errorf("KickFromGroup: find group failed, GroupID: %d, Error: %v", in.GroupId, err)
		return &group.KickFromGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 获取操作者和被踢用户的成员信息
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
		return &group.KickFromGroupResponse{Success: false, Message: "您不是该群成员，无权操作"}, nil
	}
	if targetMember == nil {
		return &group.KickFromGroupResponse{Success: false, Message: "目标用户不是该群成员"}, nil
	}

	if targetMember.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.KickFromGroupResponse{Success: false, Message: "不能将群主踢出群组"}, nil
	}
	if operatorMember.Role == int64(group.MemberRole_ROLE_ADMIN) && targetMember.Role != int64(group.MemberRole_ROLE_MEMBER) {
		return &group.KickFromGroupResponse{Success: false, Message: "权限不足，管理员只能移出普通成员"}, nil
	}
	if operatorMember.Role == int64(group.MemberRole_ROLE_MEMBER) {
		return &group.KickFromGroupResponse{Success: false, Message: "权限不足，普通成员无法移出他人"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除成员记录
	if err := tx.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Delete(&model.GroupMembers{}).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("KickFromGroup: delete member failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "移除成员失败"}, nil
	}

	// 更新群成员数量
	if err := tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).UpdateColumn("member_count", gorm.Expr("member_count - 1")).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("KickFromGroup: update member count failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "更新群成员数失败"}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("KickFromGroup: commit transaction failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "处理失败"}, nil
	}

	go l.notifyKickedUser(targetMember.UserId, &targetGroup, operatorMember)
	go l.notifyAdminsOfKickingEvent(&targetGroup, operatorMember, targetMember)

	return &group.KickFromGroupResponse{
		Success: true,
		Message: "已成功将该成员移出群组",
	}, nil
}

// notifyKickedUser 异步发送一个独立的系统通知给被踢出的用户
func (l *KickFromGroupLogic) notifyKickedUser(kickedUserID int64, groupInfo *model.Groups, operatorInfo *model.GroupMembers) {
	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	title := "您已被移出群聊"
	content := fmt.Sprintf("您已被管理员 '%s' 从群聊 '%s' 中移出。", operatorInfo.Nickname, groupInfo.Name)

	extraData := map[string]interface{}{
		"notification_type": group.NotificationType_NOTIFY_MEMBER_KICKED.String(),
		"group_id":          groupInfo.Id,
		"group_name":        groupInfo.Name,
		"operator_id":       operatorInfo.UserId,
		"operator_nickname": operatorInfo.Nickname,
	}
	extraJSON, _ := json.Marshal(extraData)

	// 只给被踢的用户发送
	err := notificationSvc.SendSystemNotification(kickedUserID, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("KickFromGroup-Notify: 发送被踢通知失败, UserID: %d, GroupID: %d, Error: %v", kickedUserID, groupInfo.Id, err)
	}
}

// notifyAdminsOfKickingEvent 异步通知管理员发生了踢人事件
func (l *KickFromGroupLogic) notifyAdminsOfKickingEvent(groupInfo *model.Groups, operatorInfo, kickedMemberInfo *model.GroupMembers) {
	// 1. 查找所有群主和管理员的ID
	var adminAndOwnerIDs []int64
	err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND role IN (?)", groupInfo.Id, []int64{int64(group.MemberRole_ROLE_OWNER), int64(group.MemberRole_ROLE_ADMIN)}).
		Pluck("user_id", &adminAndOwnerIDs).Error
	if err != nil {
		l.Logger.Errorf("KickFromGroup-AdminNotify: 查找管理员列表失败, groupID: %d, err: %v", groupInfo.Id, err)
		return
	}

	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	title := "群成员移除通知"
	content := fmt.Sprintf("在群聊 '%s' 中，管理员 '%s' 将成员 '%s' 移出。", groupInfo.Name, operatorInfo.Nickname, kickedMemberInfo.Nickname)

	extraData := map[string]interface{}{
		"notification_type": group.NotificationType_NOTIFY_MEMBER_KICKED.String(),
		"group_id":          groupInfo.Id,
		"group_name":        groupInfo.Name,
		"operator_id":       operatorInfo.UserId,
		"operator_nickname": operatorInfo.Nickname,
		"kicked_user_id":    kickedMemberInfo.UserId,
		"kicked_nickname":   kickedMemberInfo.Nickname,
	}
	extraJSON, _ := json.Marshal(extraData)

	err = notificationSvc.SendBatchNotification(adminAndOwnerIDs, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("KickFromGroup-AdminNotify: 批量发送管理员通知失败: %v", err)
	} else {
		l.Logger.Infof("KickFromGroup-AdminNotify: 已成功向群 %d 的 %d 位管理员发送成员移除通知。", groupInfo.Id, len(adminAndOwnerIDs))
	}
}
