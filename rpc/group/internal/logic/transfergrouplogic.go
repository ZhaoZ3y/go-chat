package logic

import (
	"IM/pkg/model"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferGroupLogic {
	return &TransferGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 转让群组
func (l *TransferGroupLogic) TransferGroup(in *group.TransferGroupRequest) (*group.TransferGroupResponse, error) {
	if in.GroupId == 0 || in.OwnerId == 0 || in.NewOwnerId == 0 {
		return &group.TransferGroupResponse{Success: false, Message: "参数错误"}, nil
	}
	if in.OwnerId == in.NewOwnerId {
		return &group.TransferGroupResponse{Success: false, Message: "不能将群组转让给自己"}, nil
	}

	var targetGroup model.Groups
	if err := l.svcCtx.DB.First(&targetGroup, in.GroupId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.TransferGroupResponse{Success: false, Message: "群组不存在"}, nil
		}
		l.Logger.Errorf("TransferGroup: find group failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}
	if targetGroup.OwnerId != in.OwnerId {
		return &group.TransferGroupResponse{Success: false, Message: "权限不足，只有当前群主才能转让群组"}, nil
	}

	var members []model.GroupMembers
	l.svcCtx.DB.Where("group_id = ? AND user_id IN ?", in.GroupId, []int64{in.OwnerId, in.NewOwnerId}).Find(&members)

	var oldOwnerMember, newOwnerMember *model.GroupMembers
	for i := range members {
		if members[i].UserId == in.OwnerId {
			oldOwnerMember = &members[i]
		}
		if members[i].UserId == in.NewOwnerId {
			newOwnerMember = &members[i]
		}
	}

	if oldOwnerMember == nil || newOwnerMember == nil {
		return &group.TransferGroupResponse{Success: false, Message: "新群主必须是当前群组的成员"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新群组表的 owner_id
	if err := tx.Model(&targetGroup).Update("owner_id", in.NewOwnerId).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("TransferGroup: update group owner failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "更新群主信息失败"}, nil
	}

	// 将原群主角色降级为普通成员
	if err := tx.Model(oldOwnerMember).Update("role", int64(group.MemberRole_ROLE_MEMBER)).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("TransferGroup: demote old owner role failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "更新原群主角色失败"}, nil
	}

	// 将新群主角色升级为群主
	if err := tx.Model(newOwnerMember).Update("role", int64(group.MemberRole_ROLE_OWNER)).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("TransferGroup: promote new owner role failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "更新新群主角色失败"}, nil
	}

	// 4.4 创建通知 (给新旧群主)
	notificationsToCreate := make([]*model.GroupNotification, 0, 2)
	// 通知新群主
	msgForNewOwner := fmt.Sprintf("您已被 '%s' 设置为群聊 '%s' 的新群主", oldOwnerMember.Nickname, targetGroup.Name)
	notificationsToCreate = append(notificationsToCreate, &model.GroupNotification{
		Type:         int64(group.NotificationType_NOTIFY_GROUP_TRANSFERRED),
		GroupId:      in.GroupId,
		OperatorId:   in.OwnerId,
		TargetUserId: in.NewOwnerId,
		Message:      msgForNewOwner,
	})
	// 通知旧群主
	msgForOldOwner := fmt.Sprintf("您已成功将群聊 '%s' 转让给 '%s'", targetGroup.Name, newOwnerMember.Nickname)
	notificationsToCreate = append(notificationsToCreate, &model.GroupNotification{
		Type:         int64(group.NotificationType_NOTIFY_GROUP_TRANSFERRED),
		GroupId:      in.GroupId,
		OperatorId:   in.OwnerId,
		TargetUserId: in.OwnerId,
		Message:      msgForOldOwner,
	})

	if err := tx.Create(notificationsToCreate).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("TransferGroup: create notifications failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "创建转让通知失败"}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("TransferGroup: commit transaction failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "处理失败"}, nil
	}

	//TODO：后续实现消息队列，类型为 NOTIFY_GROUP_TRANSFERRED，异步通知所有群成员该群组已转让。

	// 6. 返回成功响应
	return &group.TransferGroupResponse{
		Success: true,
		Message: "群组转让成功",
	}, nil
}
