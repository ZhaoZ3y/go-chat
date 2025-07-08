package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/chat_service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"

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
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		logx.Errorf("DismissGroup: begin transaction failed, error: %v", tx.Error)
		return nil, tx.Error
	}
	// 确保在函数退出时，如果事务未提交，则回滚
	defer tx.Rollback()

	var groupModel model.Groups
	if err := tx.Where("id = ?", in.GroupId).First(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.DismissGroupResponse{Success: false, Message: "群组不存在"}, nil
		}
		logx.Errorf("DismissGroup: find group failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 验证操作者是否为群主
	if groupModel.OwnerId != in.OwnerId {
		return &group.DismissGroupResponse{Success: false, Message: "无权解散群组，操作者非群主"}, nil
	}

	// 在删除成员之前，获取所有成员的ID用于后续通知
	var memberIDs []int64
	if err := tx.Model(&model.GroupMembers{}).Where("group_id = ?", in.GroupId).Pluck("user_id", &memberIDs).Error; err != nil {
		logx.Errorf("DismissGroup: find group members for notification failed, group_id: %d, error: %v", in.GroupId, err)
	}

	// 删除所有群组成员
	if err := tx.Where("group_id = ?", in.GroupId).Delete(&model.GroupMembers{}).Error; err != nil {
		logx.Errorf("DismissGroup: delete group members failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "删除群组成员失败"}, nil
	}

	// 删除所有相关的入群申请
	if err := tx.Where("to_group_id = ?", in.GroupId).Delete(&model.JoinGroupApplications{}).Error; err != nil {
		logx.Errorf("DismissGroup: delete join applications failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "清理入群申请失败"}, nil
	}

	// 删除群组本身
	if err := tx.Delete(&groupModel).Error; err != nil {
		logx.Errorf("DismissGroup: delete group record failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "删除群组记录失败"}, nil
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logx.Errorf("DismissGroup: commit transaction failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "解散群组事务提交失败"}, nil
	}

	// 【新增】事务提交成功后，异步发送群组解散通知
	if len(memberIDs) > 0 {
		go l.sendDismissalNotification(memberIDs, &groupModel)
	}

	return &group.DismissGroupResponse{
		Success: true,
		Message: "群组已成功解散",
	}, nil
}

func (l *DismissGroupLogic) sendDismissalNotification(userIDs []int64, groupInfo *model.Groups) {
	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	title := "群组通知"
	content := fmt.Sprintf("您所在的群组 '%s' 已被群主解散。", groupInfo.Name)

	extraData := map[string]interface{}{
		"type":        group.NotificationType_NOTIFY_GROUP_DISMISSED.String(),
		"group_id":    groupInfo.Id,
		"group_name":  groupInfo.Name,
		"operator_id": groupInfo.OwnerId,
	}
	extraJSON, err := json.Marshal(extraData)
	if err != nil {
		l.Logger.Errorf("sendDismissalNotification: failed to marshal extra data for group %d: %v", groupInfo.Id, err)
		return
	}

	// 调用批量发送接口
	err = notificationSvc.SendBatchNotification(userIDs, title, content, string(extraJSON))
	if err != nil {
		l.Logger.Errorf("sendDismissalNotification: failed to send notifications for dismissed group %d: %v", groupInfo.Id, err)
	} else {
		l.Logger.Infof("sendDismissalNotification: successfully queued dismissal notification for group %d to %d members.", groupInfo.Id, len(userIDs))
	}
}
