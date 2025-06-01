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
	// 检查群组是否存在
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND status = 1", in.GroupId).First(&groupInfo).Error; err != nil {
		return &group.JoinGroupResponse{
			Success: false,
			Message: "群组不存在",
		}, nil
	}

	// 检查是否已经是成员
	var existMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&existMember).Error; err == nil {
		return &group.JoinGroupResponse{
			Success: false,
			Message: "已经是群成员",
		}, nil
	}

	// 检查群成员数量限制
	if groupInfo.MemberCount >= groupInfo.MaxMemberCount {
		return &group.JoinGroupResponse{
			Success: false,
			Message: "群组成员已满",
		}, nil
	}

	// 获取用户信息
	var userInfo model.User
	if err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&userInfo).Error; err != nil {
		return &group.JoinGroupResponse{
			Success: false,
			Message: "用户不存在",
		}, nil
	}

	// 添加新成员
	tx := l.svcCtx.DB.Begin()
	member := &model.GroupMembers{
		GroupId:  in.GroupId,
		UserId:   in.UserId,
		Role:     3, // 普通成员
		Status:   1,
		JoinTime: time.Now().Unix(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return &group.JoinGroupResponse{
			Success: false,
			Message: "加入群组失败",
		}, nil
	}

	// 更新群成员数量
	tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
		Update("member_count", gorm.Expr("member_count + 1"))

	tx.Commit()

	// 发送站外通知给群主和管理员
	adminNotifyEvent := &notify.NotifyEvent{
		Type:      notify.NotifyTypeJoinRequest,
		GroupID:   in.GroupId,
		GroupName: groupInfo.Name,
		Data: &notify.JoinRequestData{
			UserID:   in.UserId,
			Username: userInfo.Username,
			Reason:   in.Reason,
		},
	}

	if err := l.svcCtx.NotifyService.SendNotifyToAdmins(adminNotifyEvent); err != nil {
		logx.Errorf("发送加入群聊通知失败: %v", err)
	}

	// 发送群内消息通知 - 所有群成员都能看到（包括新加入的成员）
	groupMessageEvent := &notify.NotifyEvent{
		Type:      notify.NotifyTypeJoinRequest, // 这里可以用新的类型来区分
		GroupID:   in.GroupId,
		GroupName: groupInfo.Name,
		Data: &notify.JoinRequestData{
			UserID:   in.UserId,
			Username: userInfo.Username,
			Reason:   in.Reason,
		},
	}

	if err := l.svcCtx.NotifyService.SendGroupMessage(groupMessageEvent); err != nil {
		logx.Errorf("发送群内加入通知失败: %v", err)
	}

	return &group.JoinGroupResponse{
		Success: true,
		Message: "加入群组成功",
	}, nil
}
