package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickFromGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickFromGroupLogic {
	return &KickFromGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 踢出群组
func (l *KickFromGroupLogic) KickFromGroup(in *group.KickFromGroupRequest) (*group.KickFromGroupResponse, error) {
	if in.GroupId == 0 || in.OperatorId == 0 || in.UserId == 0 {
		return &group.KickFromGroupResponse{Success: false, Message: "参数错误"}, nil
	}

	if in.OperatorId == in.UserId {
		return &group.KickFromGroupResponse{Success: false, Message: "不能将自己踢出群组"}, nil
	}

	var targetGroup model.Groups
	if err := l.svcCtx.DB.First(&targetGroup, in.GroupId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.KickFromGroupResponse{Success: false, Message: "群组不存在"}, nil
		}
		l.Logger.Errorf("KickFromGroup: find group failed, GroupID: %d, Error: %v", in.GroupId, err)
		return &group.KickFromGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 3. 获取操作者和被踢用户的成员信息，用于权限判断
	var members []model.GroupMembers
	l.svcCtx.DB.Where("group_id = ? AND user_id IN ?", in.GroupId, []int64{in.OperatorId, in.UserId}).Find(&members)

	var operatorMember, targetMember *model.GroupMembers
	for i := range members {
		if members[i].UserId == in.OperatorId {
			operatorMember = &members[i]
		}
		if members[i].UserId == in.UserId {
			targetMember = &members[i]
		}
	}

	if operatorMember == nil {
		return &group.KickFromGroupResponse{Success: false, Message: "您不是该群成员，无权操作"}, nil
	}
	if targetMember == nil {
		return &group.KickFromGroupResponse{Success: false, Message: "目标用户不是该群成员"}, nil
	}

	// 规则：群主可以踢任何人；管理员可以踢普通成员；不能踢比自己等级高或同级的人。
	if operatorMember.Role <= targetMember.Role && operatorMember.Role != int64(group.MemberRole_ROLE_OWNER) {
		return &group.KickFromGroupResponse{Success: false, Message: "权限不足，无法踢出该成员"}, nil
	}
	if targetMember.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.KickFromGroupResponse{Success: false, Message: "不能将群主踢出群组"}, nil
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
		l.Logger.Errorf("KickFromGroup: delete member failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "移除成员失败"}, nil
	}

	// 更新群成员数量
	if err := tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).UpdateColumn("member_count", gorm.Expr("member_count - 1")).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("KickFromGroup: update member count failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "更新群成员数失败"}, nil
	}

	// 创建一条通知给被踢出的用户
	notificationMessage := fmt.Sprintf("您已被管理员'%s'移出群聊 '%s'", operatorMember.Nickname, targetGroup.Name)
	notification := &model.GroupNotification{
		Type:         int64(group.NotificationType_NOTIFY_MEMBER_KICKED),
		GroupId:      in.GroupId,
		OperatorId:   in.OperatorId,
		TargetUserId: in.UserId, // 通知发给被踢的用户
		Message:      notificationMessage,
	}
	if err := tx.Create(notification).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("KickFromGroup: create notification failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "创建移除通知失败"}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("KickFromGroup: commit transaction failed: %v", err)
		return &group.KickFromGroupResponse{Success: false, Message: "处理失败"}, nil
	}

	//TODO：异步发送给群主和管理员，该用户被该高级成员踢出群聊， 类型：NOTIFY_MEMBER_KICKED

	return &group.KickFromGroupResponse{
		Success: true,
		Message: "已成功将该成员移出群组",
	}, nil
}
