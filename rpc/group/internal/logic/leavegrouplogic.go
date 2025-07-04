package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
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

	// 获取要退出的成员信息，用于后续判断
	var leavingMember model.GroupMembers
	err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&leavingMember).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.LeaveGroupResponse{Success: false, Message: "您不是该群组成员"}, nil
		}
		l.Logger.Errorf("LeaveGroup: find member failed, UserID: %d, GroupID: %d, Error: %v", in.UserId, in.GroupId, err)
		return &group.LeaveGroupResponse{Success: false, Message: "查询成员信息失败"}, nil
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

	// 查找所有群主和管理员
	var adminsToNotify []model.GroupMembers
	roles := []int64{int64(group.MemberRole_ROLE_OWNER), int64(group.MemberRole_ROLE_ADMIN)}
	tx.Where("group_id = ? AND role IN ?", in.GroupId, roles).Find(&adminsToNotify)

	notificationsToCreate := make([]*model.GroupNotification, 0)
	notificationMessage := fmt.Sprintf("成员'%s'已退出群聊", leavingMember.Nickname)

	for _, admin := range adminsToNotify {
		notificationsToCreate = append(notificationsToCreate, &model.GroupNotification{
			Type:         int64(group.NotificationType_NOTIFY_MEMBER_LEAVE),
			GroupId:      in.GroupId,
			OperatorId:   in.UserId,    // 操作者是退出者本人
			TargetUserId: admin.UserId, // 通知发给管理员
			Message:      notificationMessage,
		})
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("LeaveGroup: commit transaction failed: %v", err)
		return &group.LeaveGroupResponse{Success: false, Message: "处理失败"}, nil
	}

	// TODO： 在这里可以添加逻辑，将退出通知 (NOTIFY_MEMBER_LEAVE) 推送给所有管理员和群主。

	return &group.LeaveGroupResponse{
		Success: true,
		Message: "已成功退出群组",
	}, nil
}
