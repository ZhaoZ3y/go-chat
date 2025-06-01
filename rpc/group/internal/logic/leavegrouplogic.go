package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq/notify"
	"context"
	"gorm.io/gorm"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

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
	// 检查是否为群主
	var member model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&member).Error; err != nil {
		return &group.LeaveGroupResponse{
			Success: false,
			Message: "不是群成员",
		}, nil
	}

	if member.Role == 1 {
		return &group.LeaveGroupResponse{
			Success: false,
			Message: "群主不能退出群组，请先转让群组",
		}, nil
	}

	// 获取用户信息和群组信息
	var userInfo model.User
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&userInfo).Error; err != nil {
		return &group.LeaveGroupResponse{
			Success: false,
			Message: "用户不存在",
		}, nil
	}
	if err := l.svcCtx.DB.Where("id = ?", in.GroupId).First(&groupInfo).Error; err != nil {
		return &group.LeaveGroupResponse{
			Success: false,
			Message: "群组不存在",
		}, nil
	}

	tx := l.svcCtx.DB.Begin()
	// 删除成员记录
	if err := tx.Delete(&member).Error; err != nil {
		tx.Rollback()
		return &group.LeaveGroupResponse{
			Success: false,
			Message: "退出群组失败",
		}, nil
	}

	// 更新群成员数量
	tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
		Update("member_count", gorm.Expr("member_count - 1"))

	tx.Commit()

	// 发送通知给群主和管理员
	notifyEvent := &notify.NotifyEvent{
		Type:      notify.NotifyTypeLeaveGroup,
		GroupID:   in.GroupId,
		GroupName: groupInfo.Name,
		Data: &notify.LeaveGroupData{
			UserID:   in.UserId,
			Username: userInfo.Username,
		},
	}

	if err := l.svcCtx.NotifyService.SendNotifyToAdmins(notifyEvent); err != nil {
		logx.Errorf("发送退出群聊通知失败: %v", err)
	}

	return &group.LeaveGroupResponse{
		Success: true,
		Message: "退出群组成功",
	}, nil
}
