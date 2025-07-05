package logic

import (
	"IM/pkg/model"
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

const redisGroupMuteKey = "group:mute:%d:%d"

func NewMuteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteMemberLogic {
	return &MuteMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

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

func (l *MuteMemberLogic) MuteMember(in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	if in.GroupId == 0 || in.OperatorId == 0 || in.UserId == 0 {
		return &group.MuteMemberResponse{Success: false, Message: "参数错误"}, nil
	}
	if in.OperatorId == in.UserId {
		return &group.MuteMemberResponse{Success: false, Message: "不能对自己进行操作"}, nil
	}
	if !in.IsUnmute && in.Duration <= 0 {
		return &group.MuteMemberResponse{Success: false, Message: "禁言时长必须大于0"}, nil
	}

	// 查询操作员和目标成员权限
	var members []model.GroupMembers
	l.svcCtx.DB.Where("group_id = ? AND user_id IN ?", in.GroupId, []int64{in.OperatorId, in.UserId}).Find(&members)

	var operator, target *model.GroupMembers
	for i := range members {
		if members[i].UserId == in.OperatorId {
			operator = &members[i]
		}
		if members[i].UserId == in.UserId {
			target = &members[i]
		}
	}
	if operator == nil {
		return &group.MuteMemberResponse{Success: false, Message: "您不是该群成员，无权操作"}, nil
	}
	if target == nil {
		return &group.MuteMemberResponse{Success: false, Message: "目标用户不是该群成员"}, nil
	}

	isAdmin := operator.Role == int64(group.MemberRole_ROLE_ADMIN)
	isOwner := operator.Role == int64(group.MemberRole_ROLE_OWNER)
	if !isAdmin && !isOwner {
		return &group.MuteMemberResponse{Success: false, Message: "权限不足"}, nil
	}
	if target.Role == int64(group.MemberRole_ROLE_OWNER) {
		return &group.MuteMemberResponse{Success: false, Message: "不能禁言群主"}, nil
	}
	if isAdmin && target.Role == int64(group.MemberRole_ROLE_ADMIN) {
		return &group.MuteMemberResponse{Success: false, Message: "管理员不能禁言管理员"}, nil
	}

	muteKey := fmt.Sprintf(redisGroupMuteKey, in.GroupId, in.UserId)
	var responseMessage, notificationMessage string

	if in.IsUnmute {
		_, err := l.svcCtx.Redis.Del(l.ctx, muteKey).Result()
		if err != nil {
			l.Logger.Errorf("MuteMember: redis del failed for key %s: %v", muteKey, err)
		}
		l.svcCtx.DB.Model(&model.GroupMembers{}).
			Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
			Update("status", group.MemberStatus_MEMBER_STATUS_NORMAL)
		l.svcCtx.MuteScheduler.Remove(in.GroupId, in.UserId)
		responseMessage = "已成功解除该成员的禁言"
		notificationMessage = fmt.Sprintf("您已被管理员'%s'解除禁言", operator.Nickname)
	} else {
		expire := time.Duration(in.Duration) * time.Second
		until := time.Now().Add(expire).Unix()
		err := l.svcCtx.Redis.Set(l.ctx, muteKey, strconv.FormatInt(until, 10), expire).Err()
		if err != nil {
			l.Logger.Errorf("MuteMember: redis set failed for key %s: %v", muteKey, err)
			return &group.MuteMemberResponse{Success: false, Message: "设置禁言失败"}, nil
		}
		l.svcCtx.DB.Model(&model.GroupMembers{}).
			Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
			Update("status", group.MemberStatus_MEMBER_STATUS_MUTED)
		l.svcCtx.MuteScheduler.Register(in.GroupId, in.UserId, expire, func() {
			svc.SyncUnmuteStatus(l.svcCtx, in.GroupId, in.UserId)
		})
		durationStr := formatDuration(in.Duration)
		responseMessage = fmt.Sprintf("已成功禁言该成员，时长：%s", durationStr)
		notificationMessage = fmt.Sprintf("您已被管理员'%s'禁言，时长：%s", operator.Nickname, durationStr)
	}

	_ = l.svcCtx.DB.Create(&model.GroupNotification{
		Type:         int64(group.NotificationType_NOTIFY_MEMBER_MUTED),
		GroupId:      in.GroupId,
		OperatorId:   in.OperatorId,
		TargetUserId: in.UserId,
		Message:      notificationMessage,
	})

	return &group.MuteMemberResponse{Success: true, Message: responseMessage}, nil
}
