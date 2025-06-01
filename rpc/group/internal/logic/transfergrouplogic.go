package logic

import (
	"IM/pkg/model"
	"IM/pkg/notify"
	"context"

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
	// 检查是否为群主
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND owner_id = ?", in.GroupId, in.OwnerId).First(&groupInfo).Error; err != nil {
		return &group.TransferGroupResponse{
			Success: false,
			Message: "只有群主可以转让群组",
		}, nil
	}

	// 检查新群主是否为群成员
	var newOwnerMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.NewOwnerId).First(&newOwnerMember).Error; err != nil {
		return &group.TransferGroupResponse{
			Success: false,
			Message: "新群主不是群成员",
		}, nil
	}

	// 获取原群主和新群主信息
	var oldOwnerInfo model.User
	var newOwnerInfo model.User
	if err := l.svcCtx.DB.Where("id = ?", in.OwnerId).First(&oldOwnerInfo).Error; err != nil {
		return &group.TransferGroupResponse{
			Success: false,
			Message: "原群主不存在",
		}, nil
	}
	if err := l.svcCtx.DB.Where("id = ?", in.NewOwnerId).First(&newOwnerInfo).Error; err != nil {
		return &group.TransferGroupResponse{
			Success: false,
			Message: "新群主不存在",
		}, nil
	}

	tx := l.svcCtx.DB.Begin()

	// 更新群组所有者
	if err := tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).Update("owner_id", in.NewOwnerId).Error; err != nil {
		tx.Rollback()
		return &group.TransferGroupResponse{
			Success: false,
			Message: "转让群组失败",
		}, nil
	}

	// 更新原群主角色为普通成员
	tx.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.OwnerId).Update("role", 3)

	// 更新新群主角色
	tx.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.NewOwnerId).Update("role", 1)

	tx.Commit()

	// 发送通知给所有群成员
	notifyEvent := &notify.NotifyEvent{
		Type:      notify.NotifyTypeTransferGroup,
		GroupID:   in.GroupId,
		GroupName: groupInfo.Name,
		Data: &notify.TransferGroupData{
			OldOwnerID:   in.OwnerId,
			OldOwnerName: oldOwnerInfo.Username,
			NewOwnerID:   in.NewOwnerId,
			NewOwnerName: newOwnerInfo.Username,
		},
	}

	if err := l.svcCtx.NotifyService.SendNotifyToAllMembers(notifyEvent); err != nil {
		logx.Errorf("发送转让群聊通知失败: %v", err)
	}

	return &group.TransferGroupResponse{
		Success: true,
		Message: "转让群组成功",
	}, nil
}
