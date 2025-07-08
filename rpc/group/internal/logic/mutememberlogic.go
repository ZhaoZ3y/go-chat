package logic

import (
	"IM/pkg/model"
	"IM/pkg/mq"
	_const "IM/pkg/utils/const"
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

// 禁言群组成员
func (l *MuteMemberLogic) MuteMember(in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	// 1. 参数校验
	if in.GroupId == 0 || in.OperatorId == 0 || in.UserId == 0 {
		return &group.MuteMemberResponse{Success: false, Message: "参数错误"}, nil
	}
	if in.OperatorId == in.UserId {
		return &group.MuteMemberResponse{Success: false, Message: "不能对自己进行操作"}, nil
	}
	if !in.IsUnmute && in.Duration <= 0 {
		return &group.MuteMemberResponse{Success: false, Message: "禁言时长必须大于0"}, nil
	}

	// 2. 获取操作员和目标成员信息
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

	// 3. 权限校验
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
	var responseMessage, groupNoticeContent string
	var groupMessage *model.Messages

	tx := l.svcCtx.DB.Begin() // 开启事务
	defer tx.Rollback()

	if in.IsUnmute {
		_, err := l.svcCtx.Redis.Del(l.ctx, muteKey).Result()
		if err != nil {
			l.Logger.Errorf("MuteMember: redis del failed for key %s: %v", muteKey, err)
		}
		tx.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Update("status", group.MemberStatus_MEMBER_STATUS_NORMAL)
		l.svcCtx.MuteScheduler.Remove(in.GroupId, in.UserId)

		responseMessage = "已成功解除该成员的禁言"
		groupNoticeContent = fmt.Sprintf("成员 '%s' 已被管理员 '%s' 解除禁言", target.Nickname, operator.Nickname)

	} else {
		expire := time.Duration(in.Duration) * time.Second
		until := time.Now().Add(expire).Unix()
		err := l.svcCtx.Redis.Set(l.ctx, muteKey, strconv.FormatInt(until, 10), expire).Err()
		if err != nil {
			l.Logger.Errorf("MuteMember: redis set failed for key %s: %v", muteKey, err)
			return &group.MuteMemberResponse{Success: false, Message: "设置禁言失败"}, nil
		}
		tx.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Update("status", group.MemberStatus_MEMBER_STATUS_MUTED)
		l.svcCtx.MuteScheduler.Register(in.GroupId, in.UserId, expire, func() {
			svc.SyncUnmuteStatus(l.svcCtx, in.GroupId, in.UserId)
		})

		durationStr := formatDuration(in.Duration)
		responseMessage = fmt.Sprintf("已成功禁言该成员，时长：%s", durationStr)
		groupNoticeContent = fmt.Sprintf("成员 '%s' 已被管理员 '%s' 禁言 %s", target.Nickname, operator.Nickname, durationStr)
	}

	groupMessage = &model.Messages{
		FromUserId: _const.System, // 系统
		GroupId:    in.GroupId,
		Content:    groupNoticeContent,
		ChatType:   _const.ChatTypeGroup, // 群聊
		Type:       _const.MsgTypeSystem,
	}
	if err := tx.Create(groupMessage).Error; err != nil {
		l.Logger.Errorf("MuteMember: create group message failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("MuteMember: commit db transaction failed: %v", err)
		return &group.MuteMemberResponse{Success: false, Message: "操作失败，请重试"}, nil
	}

	if groupMessage.Id > 0 {
		go l.publishGroupSystemMessage(groupMessage.Id)
	}

	return &group.MuteMemberResponse{Success: true, Message: responseMessage}, nil
}

// publishGroupSystemMessage 异步发布群内系统消息
func (l *MuteMemberLogic) publishGroupSystemMessage(messageId int64) {
	var finalMessage model.Messages
	if err := l.svcCtx.DB.First(&finalMessage, messageId).Error; err != nil {
		l.Logger.Errorf("MuteMember-Publish: 查询群系统消息失败: %v", err)
		return
	}
	event := &mq.MessageEvent{
		Type:        mq.EventNewMessage,
		MessageID:   finalMessage.Id,
		FromUserID:  finalMessage.FromUserId,
		GroupID:     finalMessage.GroupId,
		ChatType:    finalMessage.ChatType,
		MessageType: finalMessage.Type,
		Content:     finalMessage.Content,
		CreateAt:    finalMessage.CreateAt,
	}
	if err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, event); err != nil {
		l.Logger.Errorf("MuteMember-Publish: 发布群系统消息到 Kafka 失败: %v", err)
	}
}
