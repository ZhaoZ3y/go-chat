package logic

import (
	"IM/pkg/model" // 仍然需要 model 来做权限验证
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type MuteMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// 定义 Redis 键的格式
const (
	redisGroupMuteKey = "group:mute:%d:%d"
)

func NewMuteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteMemberLogic {
	return &MuteMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MuteMember 禁言/解禁群组成员 (使用 Redis)
func (l *MuteMemberLogic) MuteMember(in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	// 参数校验
	if in.GroupId == 0 || in.OperatorId == 0 || in.UserId == 0 {
		return &group.MuteMemberResponse{Success: false, Message: "参数错误"}, nil
	}
	if in.OperatorId == in.UserId {
		return &group.MuteMemberResponse{Success: false, Message: "不能对自己进行操作"}, nil
	}
	// 禁言时长不能为负数
	if !in.IsUnmute && in.Duration <= 0 {
		return &group.MuteMemberResponse{Success: false, Message: "禁言时长必须大于0"}, nil
	}

	// 查询操作员和目标成员权限
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
		return &group.MuteMemberResponse{Success: false, Message: "您不是该群成员，无权操作"}, nil
	}
	if targetMember == nil {
		return &group.MuteMemberResponse{Success: false, Message: "目标用户不是该群成员"}, nil
	}

	isOperatorAdmin := operatorMember.Role == int64(group.MemberRole_ROLE_ADMIN)
	isOperatorOwner := operatorMember.Role == int64(group.MemberRole_ROLE_OWNER)
	if !isOperatorAdmin && !isOperatorOwner {
		return &group.MuteMemberResponse{Success: false, Message: "权限不足，只有管理员或群主才能禁言"}, nil
	}
	if targetMember.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.MuteMemberResponse{Success: false, Message: "不能禁言群主"}, nil
	}
	if isOperatorAdmin && targetMember.Role == int64(group.MemberRole_ROLE_ADMIN) {
		return &group.MuteMemberResponse{Success: false, Message: "管理员之间不能互相禁言"}, nil
	}

	muteKey := fmt.Sprintf(redisGroupMuteKey, in.GroupId, in.UserId)
	var responseMessage string
	var notificationMessage string

	if in.IsUnmute {
		// 解禁逻辑：删除 Redis 键
		_, err := l.svcCtx.Redis.Del(l.ctx, muteKey).Result()
		if err != nil {
			l.Logger.Errorf("MuteMember: redis del failed for key %s: %v", muteKey, err)
		}
		responseMessage = "已成功解除该成员的禁言"
		notificationMessage = fmt.Sprintf("您已被管理员'%s'解除禁言", operatorMember.Nickname)
	} else {
		// 禁言逻辑：设置 Redis 键和过期时间
		muteUntil := time.Now().Unix() + in.Duration
		err := l.svcCtx.Redis.Set(
			l.ctx,
			muteKey,
			strconv.FormatInt(muteUntil, 10),
			time.Duration(in.Duration)*time.Second,
		).Err()
		if err != nil {
			l.Logger.Errorf("MuteMember: redis set failed for key %s: %v", muteKey, err)
			return &group.MuteMemberResponse{Success: false, Message: "设置禁言失败"}, nil
		}
		durationStr := formatDuration(in.Duration)
		responseMessage = fmt.Sprintf("已成功禁言该成员，时长：%s", durationStr)
		notificationMessage = fmt.Sprintf("您已被管理员'%s'禁言，时长：%s", operatorMember.Nickname, durationStr)
	}

	// 创建通知（不影响核心逻辑）
	notification := &model.GroupNotification{
		Type:         int64(group.NotificationType_NOTIFY_MEMBER_MUTED),
		GroupId:      in.GroupId,
		OperatorId:   in.OperatorId,
		TargetUserId: in.UserId,
		Message:      notificationMessage,
	}
	if err := l.svcCtx.DB.Create(notification).Error; err != nil {
		l.Logger.Errorf("MuteMember: create notification failed, but mute status was set in Redis. Error: %v", err)
	}

	return &group.MuteMemberResponse{
		Success: true,
		Message: responseMessage,
	}, nil
}

// formatDuration 辅助函数 (保持不变)
func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%d分钟", seconds/60)
	}
	if seconds < 86400 {
		return fmt.Sprintf("%d小时", seconds/3600)
	}
	return fmt.Sprintf("%d天", seconds/86400)
}
