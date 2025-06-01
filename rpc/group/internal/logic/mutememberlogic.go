package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq/notify"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MuteMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMuteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteMemberLogic {
	return &MuteMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁言群组成员
func (l *MuteMemberLogic) MuteMember(in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	// 1. 检查操作者权限
	var operatorMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role IN (1,2)",
		in.GroupId, in.OperatorId).First(&operatorMember).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "无权限操作"}, nil
	}

	// 2. 检查目标用户
	var targetMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&targetMember).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "用户不在群组中"}, nil
	}
	if targetMember.Role == 1 {
		return &group.MuteMemberResponse{Success: false, Message: "不能禁言群主"}, nil
	}
	if operatorMember.Role == 2 && targetMember.Role == 2 {
		return &group.MuteMemberResponse{Success: false, Message: "管理员不能禁言其他管理员"}, nil
	}

	// 3. 获取用户和群信息
	var (
		operatorInfo   model.User
		targetUserInfo model.User
		groupInfo      model.Groups
	)
	if err := l.svcCtx.DB.Where("id = ?", in.OperatorId).First(&operatorInfo).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "操作者不存在"}, nil
	}
	if err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&targetUserInfo).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "目标用户不存在"}, nil
	}
	if err := l.svcCtx.DB.Where("id = ?", in.GroupId).First(&groupInfo).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "群组不存在"}, nil
	}

	// 4. 构造 Redis Key
	muteKey := fmt.Sprintf("mute:group:%d:user:%d", in.GroupId, in.UserId)

	// 5. 事务处理（数据库 + Redis）
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		ctx := context.Background()

		if in.Duration > 0 {
			// 禁言逻辑
			muteUntil := time.Now().Add(time.Duration(in.Duration) * time.Second)
			if err := tx.Model(&model.GroupMembers{}).
				Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
				Update("mute_until", muteUntil).Error; err != nil {
				return err
			}

			muteData := map[string]interface{}{
				"operator_id": in.OperatorId,
				"username":    targetUserInfo.Username,
				"duration":    in.Duration,
			}
			data, _ := json.Marshal(muteData)

			if err := l.svcCtx.Redis.Set(ctx, muteKey, data, time.Duration(in.Duration)*time.Second).Err(); err != nil {
				return err
			}
		} else {
			// 手动解除禁言逻辑
			if err := tx.Model(&model.GroupMembers{}).
				Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
				Update("mute_until", time.Time{}).Error; err != nil {
				return err
			}
			_ = l.svcCtx.Redis.Del(ctx, muteKey).Err()
		}

		// 通知推送
		notifyEvent := &notify.NotifyEvent{
			Type:      notify.NotifyTypeMuteMember,
			GroupID:   in.GroupId,
			GroupName: groupInfo.Name,
			Data: &notify.MuteMemberData{
				OperatorID:   in.OperatorId,
				OperatorName: operatorInfo.Username,
				UserID:       in.UserId,
				Username:     targetUserInfo.Username,
				Duration:     in.Duration,
			},
		}

		// 发送群内消息 - 通过WebSocket发送给所有群成员
		if err := l.svcCtx.NotifyService.SendGroupMessage(notifyEvent); err != nil {
			logx.Errorf("发送群内通知失败: %v", err)
			// 继续执行，不因通知失败而回滚事务
		}

		return nil
	})

	if err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "操作失败"}, nil
	}

	// 成功返回
	msg := "禁言成功"
	if in.Duration == 0 {
		msg = "解除禁言成功"
	}

	return &group.MuteMemberResponse{
		Success: true,
		Message: msg,
	}, nil
}
