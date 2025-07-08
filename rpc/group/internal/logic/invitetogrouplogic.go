package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"IM/pkg/utils/chat_service"
	_const "IM/pkg/utils/const"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
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

	// 校验被邀请用户是否存在
	var existingUsers []model.User
	if err := tx.Model(&model.User{}).Where("id IN ?", in.UserIds).Find(&existingUsers).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("InviteToGroup: 查询用户失败: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "查询被邀请用户失败"}, nil
	}

	existingUserMap := make(map[int64]*model.User, len(existingUsers))
	for i := range existingUsers {
		existingUserMap[existingUsers[i].Id] = &existingUsers[i]
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
		return &group.InviteToGroupResponse{Success: false, Message: msg}, nil
	}

	// 查找已在群中的人
	var existingMembers []model.GroupMembers
	tx.Where("group_id = ? AND user_id IN ?", in.GroupId, in.UserIds).Find(&existingMembers)
	existingMemberMap := make(map[int64]struct{}, len(existingMembers))
	for _, m := range existingMembers {
		existingMemberMap[m.UserId] = struct{}{}
	}

	// 查找已有待处理入群申请的用户
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

	var groupNotificationMessage *model.Messages
	var newApplicationsCreated bool // 标记是否有新的申请被创建

	if isAdmin {
		// 管理员或群主，直接加入群组
		membersToCreate := make([]*model.GroupMembers, 0)
		newMemberNicknames := make([]string, 0)

		for _, uid := range in.UserIds {
			if _, exists := existingMemberMap[uid]; exists {
				failedUserIDs = append(failedUserIDs, uid)
				continue
			}
			user, _ := existingUserMap[uid]
			membersToCreate = append(membersToCreate, &model.GroupMembers{
				GroupId: in.GroupId, UserId: uid, Nickname: user.Nickname,
				Role: int64(group.MemberRole_ROLE_MEMBER), Status: int64(group.MemberStatus_MEMBER_STATUS_NORMAL), JoinTime: time.Now().Unix(),
			})
			newMemberNicknames = append(newMemberNicknames, user.Nickname)
		}

		if len(membersToCreate) > 0 {
			if err := tx.Create(&membersToCreate).Error; err != nil {
				tx.Rollback()
				return &group.InviteToGroupResponse{Success: false, Message: "添加新成员失败"}, nil
			}
			if err := tx.Model(&targetGroup).UpdateColumn("member_count", gorm.Expr("member_count + ?", len(membersToCreate))).Error; err != nil {
				tx.Rollback()
				return &group.InviteToGroupResponse{Success: false, Message: "更新群成员数失败"}, nil
			}

			notificationContent := fmt.Sprintf("'%s' 邀请 '%s' 加入了群聊", inviterMember.Nickname, strings.Join(newMemberNicknames, "、"))
			groupNotificationMessage = &model.Messages{
				FromUserId: _const.System,
				GroupId:    in.GroupId,
				Content:    notificationContent,
				ChatType:   _const.ChatTypeGroup,
				Type:       _const.MsgTypeSystem,
			}
			if err := tx.Create(groupNotificationMessage).Error; err != nil {
				tx.Rollback()
				logx.Errorf("InviteToGroup: 创建群内通知消息失败: %v", err)
			}
		}
		message = fmt.Sprintf("成功邀请 %d 人加入群组，%d 人失败或已在群中。", len(membersToCreate), len(failedUserIDs))

	} else {
		// 普通成员，创建入群申请
		applicationsToCreate := make([]*model.JoinGroupApplications, 0)
		for _, uid := range in.UserIds {
			if _, exists := existingMemberMap[uid]; exists {
				failedUserIDs = append(failedUserIDs, uid)
				continue
			}
			if _, exists := pendingApplicationMap[uid]; exists {
				failedUserIDs = append(failedUserIDs, uid)
				continue
			}
			applicationsToCreate = append(applicationsToCreate, &model.JoinGroupApplications{
				FromUserId: uid, ToGroupId: in.GroupId,
				Reason:    fmt.Sprintf("由成员 '%s' 邀请加入", inviterMember.Nickname),
				InviterId: in.InviterId, Status: int8(group.ApplicationStatus_PENDING),
			})
		}

		if len(applicationsToCreate) > 0 {
			if err := tx.Create(&applicationsToCreate).Error; err != nil {
				tx.Rollback()
				return &group.InviteToGroupResponse{Success: false, Message: "创建入群申请失败"}, nil
			}
			newApplicationsCreated = true // 标记成功创建了新申请
		}
		message = fmt.Sprintf("已为 %d 人发送入群申请，%d 人失败或已申请/在群中。", len(applicationsToCreate), len(failedUserIDs))
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("InviteToGroup: 提交事务失败: %v", err)
		return &group.InviteToGroupResponse{Success: false, Message: "处理邀请失败"}, nil
	}

	// 如果是管理员直接拉人，则在群聊内发系统消息
	if groupNotificationMessage != nil {
		go l.publishGroupSystemMessage(groupNotificationMessage.Id)
	}

	// 如果是普通成员创建了新的入群申请，则通知所有管理员
	if newApplicationsCreated {
		go l.notifyAdminsOfNewApplication(&targetGroup, &inviterMember)
	}

	return &group.InviteToGroupResponse{
		Success:       true,
		Message:       message,
		FailedUserIds: failedUserIDs,
	}, nil
}

// publishGroupSystemMessage 用于异步发布群内系统消息
func (l *InviteToGroupLogic) publishGroupSystemMessage(messageId int64) {
	var finalMessage model.Messages
	if err := l.svcCtx.DB.First(&finalMessage, messageId).Error; err != nil {
		l.Logger.Errorf("发布群系统消息前查询消息体失败: %v", err)
		return
	}

	event := &mq.MessageEvent{
		Type:        mq.EventNewMessage,
		MessageID:   finalMessage.Id,
		FromUserID:  finalMessage.FromUserId,
		GroupID:     finalMessage.GroupId,
		ChatType:    finalMessage.ChatType,
		MessageType: finalMessage.Type,
		Content:     finalMessage.Content,
		CreateAt:    finalMessage.CreateAt,
	}

	if err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, event); err != nil {
		l.Logger.Errorf("发布群系统消息到 Kafka 失败: %v", err)
	}
}

// notifyAdminsOfNewApplication 用于异步通知管理员有新的入群申请
func (l *InviteToGroupLogic) notifyAdminsOfNewApplication(groupInfo *model.Groups, inviterInfo *model.GroupMembers) {
	// 1. 查找所有群主和管理员的ID
	var adminAndOwnerIDs []int64
	err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND role IN (?)", groupInfo.Id, []int64{int64(group.MemberRole_ROLE_OWNER), int64(group.MemberRole_ROLE_ADMIN)}).
		Pluck("user_id", &adminAndOwnerIDs).Error

	if err != nil {
		l.Logger.Errorf("通知管理员失败：查找管理员列表时出错, groupID: %d, err: %v", groupInfo.Id, err)
		return
	}

	// 2. 如果没有管理员，则无需通知
	if len(adminAndOwnerIDs) == 0 {
		return
	}

	// 3. 准备并发送通知
	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	title := "新的入群申请"
	content := fmt.Sprintf("成员 '%s' 邀请了新成员加入群聊 '%s'，请及时处理。", inviterInfo.Nickname, groupInfo.Name)

	extraData := map[string]interface{}{
		"notification_type": group.NotificationType_NOTIFY_MEMBER_APPLY_JOIN.String(),
		"group_id":          groupInfo.Id,
		"group_name":        groupInfo.Name,
		"inviter_id":        inviterInfo.UserId,
		"inviter_nickname":  inviterInfo.Nickname,
	}
	extraJSON, _ := json.Marshal(extraData)

	// 向所有管理员和群主批量发送这条独立的系统通知
	err = notificationSvc.SendBatchNotification(adminAndOwnerIDs, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("向管理员批量发送新申请通知失败: %v", err)
	} else {
		l.Logger.Infof("已成功向群 %d 的 %d 位管理员发送新申请通知。", groupInfo.Id, len(adminAndOwnerIDs))
	}
}
