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

type LeaveGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupLogic {
	return &LeaveGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 退出群组
func (l *LeaveGroupLogic) LeaveGroup(in *group.LeaveGroupRequest) (*group.LeaveGroupResponse, error) {
	if in.GroupId == 0 || in.UserId == 0 {
		return &group.LeaveGroupResponse{Success: false, Message: "参数错误"}, nil
	}

	// 获取要退出的成员信息和群组信息
	var leavingMember model.GroupMembers
	err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&leavingMember).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.LeaveGroupResponse{Success: false, Message: "您不是该群组成员"}, nil
		}
		l.Logger.Errorf("LeaveGroup: find member failed, UserID: %d, GroupID: %d, Error: %v", in.UserId, in.GroupId, err)
		return &group.LeaveGroupResponse{Success: false, Message: "查询成员信息失败"}, nil
	}

	var targetGroup model.Groups
	if err := l.svcCtx.DB.First(&targetGroup, in.GroupId).Error; err != nil {
		l.Logger.Errorf("LeaveGroup: find group failed, GroupID: %d, Error: %v", in.GroupId, err)
		return &group.LeaveGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	if leavingMember.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.LeaveGroupResponse{Success: false, Message: "群主不能直接退出群组，请先转让群主或解散群组"}, nil
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
		l.Logger.Errorf("LeaveGroup: delete member failed: %v", err)
		return &group.LeaveGroupResponse{Success: false, Message: "退出群组失败"}, nil
	}

	// 更新群成员数量
	if err := tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).UpdateColumn("member_count", gorm.Expr("member_count - 1")).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("LeaveGroup: update member count failed: %v", err)
		return &group.LeaveGroupResponse{Success: false, Message: "更新群成员数失败"}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("LeaveGroup: commit transaction failed: %v", err)
		return &group.LeaveGroupResponse{Success: false, Message: "处理失败"}, nil
	}

	go l.notifyAdminsOfLeaveEvent(&targetGroup, &leavingMember)

	return &group.LeaveGroupResponse{
		Success: true,
		Message: "已成功退出群组",
	}, nil
}

// notifyAdminsOfLeaveEvent 异步通知管理员有成员退出了群组
func (l *LeaveGroupLogic) notifyAdminsOfLeaveEvent(groupInfo *model.Groups, leavingMemberInfo *model.GroupMembers) {
	var adminAndOwnerIDs []int64
	err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND role IN (?)", groupInfo.Id, []int64{int64(group.MemberRole_ROLE_OWNER), int64(group.MemberRole_ROLE_ADMIN)}).
		Pluck("user_id", &adminAndOwnerIDs).Error
	if err != nil {
		l.Logger.Errorf("LeaveGroup-Notify: 查找管理员列表失败, groupID: %d, err: %v", groupInfo.Id, err)
		return
	}

	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	title := "群成员退出通知"
	content := fmt.Sprintf("成员 '%s' 已主动退出群聊 '%s'。", leavingMemberInfo.Nickname, groupInfo.Name)

	extraData := map[string]interface{}{
		"notification_type": group.NotificationType_NOTIFY_MEMBER_LEAVE.String(), // "NOTIFY_MEMBER_LEAVE"
		"group_id":          groupInfo.Id,
		"group_name":        groupInfo.Name,
		"leaving_user_id":   leavingMemberInfo.UserId,
		"leaving_nickname":  leavingMemberInfo.Nickname,
	}
	extraJSON, _ := json.Marshal(extraData)

	err = notificationSvc.SendBatchNotification(adminAndOwnerIDs, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("LeaveGroup-Notify: 批量发送管理员通知失败: %v", err)
	} else {
		l.Logger.Infof("LeaveGroup-Notify: 已成功向群 %d 的 %d 位管理员发送成员退出通知。", groupInfo.Id, len(adminAndOwnerIDs))
	}
}
