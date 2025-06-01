package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq/notify"
	"context"
	"gorm.io/gorm"
	"time"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

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
	// 检查邀请者权限
	var inviterMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role IN (1,2)",
		in.GroupId, in.InviterId).First(&inviterMember).Error; err != nil {
		return &group.InviteToGroupResponse{
			Success: false,
			Message: "无权限邀请",
		}, nil
	}

	// 获取邀请者信息和群组信息
	var inviterInfo model.User
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ?", in.InviterId).First(&inviterInfo).Error; err != nil {
		return &group.InviteToGroupResponse{
			Success: false,
			Message: "邀请者不存在",
		}, nil
	}
	if err := l.svcCtx.DB.Where("id = ?", in.GroupId).First(&groupInfo).Error; err != nil {
		return &group.InviteToGroupResponse{
			Success: false,
			Message: "群组不存在",
		}, nil
	}

	var failedUserIds []int64
	var successUserIds []int64
	var successUsernames []string
	tx := l.svcCtx.DB.Begin()

	for _, userId := range in.UserIds {
		// 检查是否已经是成员
		var existMember model.GroupMembers
		if err := tx.Where("group_id = ? AND user_id = ?", in.GroupId, userId).First(&existMember).Error; err == nil {
			failedUserIds = append(failedUserIds, userId)
			continue
		}

		// 获取被邀请用户信息
		var userInfo model.User
		if err := tx.Where("id = ?", userId).First(&userInfo).Error; err != nil {
			failedUserIds = append(failedUserIds, userId)
			continue
		}

		// 添加新成员
		member := &model.GroupMembers{
			GroupId:  in.GroupId,
			UserId:   userId,
			Role:     3,
			Status:   1,
			JoinTime: time.Now().Unix(),
		}

		if err := tx.Create(member).Error; err != nil {
			failedUserIds = append(failedUserIds, userId)
			continue
		}

		// 更新群成员数量
		tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
			Update("member_count", gorm.Expr("member_count + 1"))

		successUserIds = append(successUserIds, userId)
		successUsernames = append(successUsernames, userInfo.Username)
	}

	tx.Commit()

	if len(successUserIds) > 0 {
		// 发送站外通知给群主和管理员
		adminNotifyEvent := &notify.NotifyEvent{
			Type:      notify.NotifyTypeInviteToGroup,
			GroupID:   in.GroupId,
			GroupName: groupInfo.Name,
			Data: &notify.InviteToGroupData{
				InviterID:   in.InviterId,
				InviterName: inviterInfo.Username,
				UserIDs:     successUserIds,
				Usernames:   successUsernames,
			},
		}

		if err := l.svcCtx.NotifyService.SendNotifyToAdmins(adminNotifyEvent); err != nil {
			logx.Errorf("发送邀请加入群聊通知失败: %v", err)
		}

		// 发送群内消息通知 - 所有群成员都能看到
		groupMessageEvent := &notify.NotifyEvent{
			Type:      notify.NotifyTypeInviteToGroup,
			GroupID:   in.GroupId,
			GroupName: groupInfo.Name,
			Data: &notify.InviteToGroupData{
				InviterID:   in.InviterId,
				InviterName: inviterInfo.Username,
				UserIDs:     successUserIds,
				Usernames:   successUsernames,
			},
		}

		if err := l.svcCtx.NotifyService.SendGroupMessage(groupMessageEvent); err != nil {
			logx.Errorf("发送群内邀请通知失败: %v", err)
		}
	}

	return &group.InviteToGroupResponse{
		Success:       true,
		Message:       "邀请完成",
		FailedUserIds: failedUserIds,
	}, nil
}
