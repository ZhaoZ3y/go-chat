package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq/notify"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DismissGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDismissGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DismissGroupLogic {
	return &DismissGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解散群组
func (l *DismissGroupLogic) DismissGroup(in *group.DismissGroupRequest) (*group.DismissGroupResponse, error) {
	// 检查是否为群主
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND owner_id = ?", in.GroupId, in.OwnerId).First(&groupInfo).Error; err != nil {
		return &group.DismissGroupResponse{
			Success: false,
			Message: "只有群主可以解散群组",
		}, nil
	}

	// 获取群主信息
	var ownerInfo model.User
	if err := l.svcCtx.DB.Where("id = ?", in.OwnerId).First(&ownerInfo).Error; err != nil {
		return &group.DismissGroupResponse{
			Success: false,
			Message: "群主不存在",
		}, nil
	}

	// 先发送通知给所有群成员
	notifyEvent := &notify.NotifyEvent{
		Type:      notify.NotifyTypeDismissGroup,
		GroupID:   in.GroupId,
		GroupName: groupInfo.Name,
		Data: &notify.DismissGroupData{
			OwnerID:   in.OwnerId,
			OwnerName: ownerInfo.Username,
		},
	}

	if err := l.svcCtx.NotifyService.SendNotifyToAllMembers(notifyEvent); err != nil {
		logx.Errorf("发送解散群聊通知失败: %v", err)
	}

	tx := l.svcCtx.DB.Begin()

	// 删除所有群成员
	if err := tx.Where("group_id = ?", in.GroupId).Delete(&model.GroupMembers{}).Error; err != nil {
		tx.Rollback()
		return &group.DismissGroupResponse{
			Success: false,
			Message: "解散群组失败",
		}, nil
	}

	// 删除群组
	if err := tx.Delete(&groupInfo).Error; err != nil {
		tx.Rollback()
		return &group.DismissGroupResponse{
			Success: false,
			Message: "解散群组失败",
		}, nil
	}

	tx.Commit()
	return &group.DismissGroupResponse{
		Success: true,
		Message: "解散群组成功",
	}, nil
}
