package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	"IM/pkg/utils/chat_service"
	_const "IM/pkg/utils/const"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	defer tx.Rollback()

	// 更新群组表的 owner_id
	if err := tx.Model(&targetGroup).Update("owner_id", in.NewOwnerId).Error; err != nil {
		tx.Rollback()
		return &group.TransferGroupResponse{Success: false, Message: "更新群主信息失败"}, nil
	}

	// 将原群主角色降级为普通成员
	if err := tx.Model(oldOwnerMember).Update("role", int64(group.MemberRole_ROLE_MEMBER)).Error; err != nil {
		tx.Rollback()
		return &group.TransferGroupResponse{Success: false, Message: "更新原群主角色失败"}, nil
	}

	// 将新群主角色升级为群主
	if err := tx.Model(newOwnerMember).Update("role", int64(group.MemberRole_ROLE_OWNER)).Error; err != nil {
		tx.Rollback()
		return &group.TransferGroupResponse{Success: false, Message: "更新新群主角色失败"}, nil
	}

	// 3在事务中创建一条群内广播消息
	groupNoticeContent := fmt.Sprintf("群主已变更为 '%s'", newOwnerMember.Nickname)
	groupMessage := &model.Messages{
		FromUserId: _const.System,
		GroupId:    in.GroupId,
		Content:    groupNoticeContent,
		ChatType:   _const.ChatTypeGroup,
		Type:       _const.MsgTypeSystem,
	}
	if err := tx.Create(groupMessage).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("TransferGroup: create group message failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("TransferGroup: commit transaction failed: %v", err)
		return &group.TransferGroupResponse{Success: false, Message: "处理失败"}, nil
	}

	go l.notifyTransferParties(oldOwnerMember, newOwnerMember, &targetGroup)
	if groupMessage.Id > 0 {
		go l.publishGroupSystemMessage(groupMessage.Id)
	}

	return &group.TransferGroupResponse{
		Success: true,
		Message: "群组转让成功",
	}, nil
}

// notifyTransferParties 异步通知原群主和新群主
func (l *TransferGroupLogic) notifyTransferParties(oldOwner, newOwner *model.GroupMembers, groupInfo *model.Groups) {
	notificationSvc := chat_service.NewNotificationService(l.svcCtx.DB, l.svcCtx.Kafka)

	// 准备通知新群主的内容
	titleForNewOwner := "您已成为新群主"
	contentForNewOwner := fmt.Sprintf("您已被 '%s' 设置为群聊 '%s' 的新群主。", oldOwner.Nickname, groupInfo.Name)
	extraDataForNew := map[string]interface{}{
		"notification_type":  group.NotificationType_NOTIFY_GROUP_TRANSFERRED.String(),
		"group_id":           groupInfo.Id,
		"group_name":         groupInfo.Name,
		"old_owner_id":       oldOwner.UserId,
		"old_owner_nickname": oldOwner.Nickname,
	}
	extraJSONForNew, _ := json.Marshal(extraDataForNew)
	err := notificationSvc.SendSystemNotification(newOwner.UserId, titleForNewOwner, contentForNewOwner, string(extraJSONForNew))
	if err != nil {
		l.Logger.Errorf("TransferGroup-Notify: 发送新群主通知失败: %v", err)
	}

	// 准备通知原群主的内容
	titleForOldOwner := "群组已成功转让"
	contentForOldOwner := fmt.Sprintf("您已成功将群聊 '%s' 转让给 '%s'。", groupInfo.Name, newOwner.Nickname)
	extraDataForOld := map[string]interface{}{
		"notification_type":  group.NotificationType_NOTIFY_GROUP_TRANSFERRED.String(),
		"group_id":           groupInfo.Id,
		"group_name":         groupInfo.Name,
		"new_owner_id":       newOwner.UserId,
		"new_owner_nickname": newOwner.Nickname,
	}
	extraJSONForOld, _ := json.Marshal(extraDataForOld)
	err = notificationSvc.SendSystemNotification(oldOwner.UserId, titleForOldOwner, contentForOldOwner, string(extraJSONForOld))
	if err != nil {
		l.Logger.Errorf("TransferGroup-Notify: 发送原群主通知失败: %v", err)
	}
}

// publishGroupSystemMessage 异步发布群内广播消息
func (l *TransferGroupLogic) publishGroupSystemMessage(messageId int64) {
	var finalMessage model.Messages
	if err := l.svcCtx.DB.First(&finalMessage, messageId).Error; err != nil {
		l.Logger.Errorf("TransferGroup-Publish: 查询群系统消息失败: %v", err)
		return
	}
	event := &mq.MessageEvent{
		Type: mq.EventNewMessage, MessageID: finalMessage.Id, FromUserID: finalMessage.FromUserId,
		GroupID: finalMessage.GroupId, ChatType: finalMessage.ChatType, MessageType: finalMessage.Type,
		Content: finalMessage.Content, CreateAt: finalMessage.CreateAt,
	}
	if err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, event); err != nil {
		l.Logger.Errorf("TransferGroup-Publish: 发布群系统消息到 Kafka 失败: %v", err)
	}
}
