package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/chat_service"
	_const "IM/pkg/utils/const"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 加入群组
func (l *JoinGroupLogic) JoinGroup(in *group.JoinGroupRequest) (*group.JoinGroupResponse, error) {
	if in.GroupId == 0 || in.UserId == 0 {
		return &group.JoinGroupResponse{Success: false, Message: "参数错误：群组ID和用户ID不能为空"}, nil
	}

	var targetGroup model.Groups
	err := l.svcCtx.DB.Where("id = ? AND status = ?", in.GroupId, group.GroupStatus_GROUP_STATUS_NORMAL).First(&targetGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.JoinGroupResponse{Success: false, Message: "群组不存在或已解散"}, nil
		}
		l.Logger.Errorf("JoinGroup: find group failed, GroupID: %d, Error: %v", in.GroupId, err)
		return &group.JoinGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 检查用户是否已经是群成员
	var memberCount int64
	l.svcCtx.DB.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Count(&memberCount)
	if memberCount > 0 {
		return &group.JoinGroupResponse{Success: false, Message: "您已是该群组成员"}, nil
	}

	// 检查是否已有待处理的申请，防止重复提交
	var applicationCount int64
	l.svcCtx.DB.Model(&model.JoinGroupApplications{}).Where("to_group_id = ? AND from_user_id = ? AND status = ?", in.GroupId, in.UserId, group.ApplicationStatus_PENDING).Count(&applicationCount)
	if applicationCount > 0 {
		return &group.JoinGroupResponse{Success: false, Message: "您已提交过申请，请耐心等待管理员审核"}, nil
	}

	var applicantUser model.User
	if err := l.svcCtx.DB.First(&applicantUser, in.UserId).Error; err != nil {
		l.Logger.Errorf("JoinGroup: 获取申请人信息失败, UserID: %d, Error: %v", in.UserId, err)
		return &group.JoinGroupResponse{Success: false, Message: "获取您的用户信息失败"}, nil
	}

	application := &model.JoinGroupApplications{
		FromUserId: in.UserId,
		ToGroupId:  in.GroupId,
		Reason:     in.Reason,
		InviterId:  _const.ApplyProactive,
		Status:     int8(group.ApplicationStatus_PENDING),
	}

	if err := l.svcCtx.DB.Create(application).Error; err != nil {
		l.Logger.Errorf("JoinGroup: create application failed: %v", err)
		return &group.JoinGroupResponse{Success: false, Message: "提交申请失败，请稍后再试"}, nil
	}

	go l.notifyAdminsOfNewApplication(&targetGroup, &applicantUser)

	return &group.JoinGroupResponse{
		Success: true,
		Message: "入群申请已提交，请等待管理员审核",
	}, nil
}

// notifyAdminsOfNewApplication 用于异步通知管理员有新的入群申请
func (l *JoinGroupLogic) notifyAdminsOfNewApplication(groupInfo *model.Groups, applicantInfo *model.User) {
	// 1. 查找所有群主和管理员的ID
	var adminAndOwnerIDs []int64
	err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND role IN (?)", groupInfo.Id, []int64{int64(group.MemberRole_ROLE_OWNER), int64(group.MemberRole_ROLE_ADMIN)}).
		Pluck("user_id", &adminAndOwnerIDs).Error

	if err != nil {
		l.Logger.Errorf("JoinGroup-Notify: 查找管理员列表失败, groupID: %d, err: %v", groupInfo.Id, err)
		return
	}

	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	title := "新的入群申请"
	content := fmt.Sprintf("用户 '%s' 申请加入群聊 '%s'，请及时处理。", applicantInfo.Nickname, groupInfo.Name)

	extraData := map[string]interface{}{
		"notification_type":  group.NotificationType_NOTIFY_MEMBER_APPLY_JOIN.String(),
		"group_id":           groupInfo.Id,
		"group_name":         groupInfo.Name,
		"applicant_id":       applicantInfo.Id,
		"applicant_nickname": applicantInfo.Nickname,
	}
	extraJSON, _ := json.Marshal(extraData)

	// 向所有管理员和群主批量发送这条独立的系统通知
	err = notificationSvc.SendBatchNotification(adminAndOwnerIDs, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("JoinGroup-Notify: 向管理员批量发送新申请通知失败: %v", err)
	} else {
		l.Logger.Infof("JoinGroup-Notify: 已成功向群 %d 的 %d 位管理员发送新申请通知。", groupInfo.Id, len(adminAndOwnerIDs))
	}
}
