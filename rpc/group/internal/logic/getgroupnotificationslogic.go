package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupNotificationsLogic {
	return &GetGroupNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupNotificationsLogic) GetGroupNotifications(in *group.GetGroupNotificationsRequest) (*group.GetGroupNotificationsResponse, error) {
	var notifications []model.GroupNotification
	err := l.svcCtx.DB.
		Where("target_user_id = ?", in.UserId).
		Order("create_at DESC").
		Find(&notifications).Error
	if err != nil {
		l.Logger.Errorf("查询群组通知失败: %v", err)
		return nil, err
	}

	// 收集相关 userId、groupId
	userIDSet := make(map[int64]struct{})
	groupIDSet := make(map[int64]struct{})

	for _, n := range notifications {
		userIDSet[n.OperatorId] = struct{}{}
		userIDSet[n.TargetUserId] = struct{}{}
		groupIDSet[n.GroupId] = struct{}{}
	}

	var userIds, groupIds []int64
	for id := range userIDSet {
		userIds = append(userIds, id)
	}
	for id := range groupIDSet {
		groupIds = append(groupIds, id)
	}

	// 批量查询用户信息
	var users []model.User
	if err := l.svcCtx.DB.Where("id IN ?", userIds).Find(&users).Error; err != nil {
		l.Logger.Errorf("查询用户信息失败: %v", err)
		return nil, err
	}
	userMap := make(map[int64]*model.User)
	for _, u := range users {
		uCopy := u
		userMap[u.Id] = &uCopy
	}

	// 批量查询群组信息
	var groups []model.Groups
	if err := l.svcCtx.DB.Where("id IN ?", groupIds).Find(&groups).Error; err != nil {
		l.Logger.Errorf("查询群组信息失败: %v", err)
		return nil, err
	}
	groupMap := make(map[int64]*model.Groups)
	for _, g := range groups {
		gCopy := g
		groupMap[g.Id] = &gCopy
	}

	// 构造响应
	var result []*group.GroupNotification
	for _, n := range notifications {
		operator := userMap[n.OperatorId]
		target := userMap[n.TargetUserId]
		groupInfo := groupMap[n.GroupId]

		pb := &group.GroupNotification{
			Id:                 n.Id,
			Type:               group.NotificationType(n.Type),
			GroupId:            n.GroupId,
			OperatorId:         n.OperatorId,
			TargetUserId:       n.TargetUserId,
			Message:            n.Message,
			Timestamp:          n.CreateAt,
			IsRead:             n.IsRead,
			OperatorNickname:   "",
			OperatorAvatar:     "",
			TargetUserNickname: "",
			TargetUserAvatar:   "",
			GroupName:          "",
			GroupAvatar:        "",
		}

		if operator != nil {
			pb.OperatorNickname = operator.Nickname
			pb.OperatorAvatar = operator.Avatar
		}
		if target != nil {
			pb.TargetUserNickname = target.Nickname
			pb.TargetUserAvatar = target.Avatar
		}
		if groupInfo != nil {
			pb.GroupName = groupInfo.Name
			pb.GroupAvatar = groupInfo.Avatar
		}

		result = append(result, pb)
	}

	return &group.GetGroupNotificationsResponse{
		Notifications: result,
	}, nil
}
