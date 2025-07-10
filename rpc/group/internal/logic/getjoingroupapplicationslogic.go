package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJoinGroupApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetJoinGroupApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJoinGroupApplicationsLogic {
	return &GetJoinGroupApplicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取加入群组申请列表
func (l *GetJoinGroupApplicationsLogic) GetJoinGroupApplications(in *group.GetJoinGroupApplicationsRequest) (*group.GetJoinGroupApplicationsResponse, error) {
	// 查询申请记录
	var applications []*model.JoinGroupApplications
	err := l.svcCtx.DB.WithContext(l.ctx).
		Order("create_at DESC").
		Find(&applications).Error
	if err != nil {
		l.Logger.Errorf("获取加群申请失败: %v", err)
		return nil, err
	}

	userIDSet := make(map[int64]struct{})
	groupIDSet := make(map[int64]struct{})

	for _, app := range applications {
		userIDSet[app.FromUserId] = struct{}{}
		groupIDSet[app.ToGroupId] = struct{}{}
		if app.InviterId > 0 {
			userIDSet[app.InviterId] = struct{}{}
		}
		if app.OperatorId > 0 {
			userIDSet[app.OperatorId] = struct{}{}
		}
	}

	// 转 slice
	var userIds []int64
	for id := range userIDSet {
		userIds = append(userIds, id)
	}
	var groupIds []int64
	for id := range groupIDSet {
		groupIds = append(groupIds, id)
	}

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

	// ----------------------------
	// 4. 构建响应
	// ----------------------------
	var pbApplications []*group.JoinGroupApplication
	for _, app := range applications {
		user := userMap[app.FromUserId]
		inviter := userMap[app.InviterId]
		operator := userMap[app.OperatorId]
		grouper := groupMap[app.ToGroupId]

		pbApp := &group.JoinGroupApplication{
			Id:               app.Id,
			UserId:           app.FromUserId,
			GroupId:          app.ToGroupId,
			Reason:           app.Reason,
			ApplyTime:        app.CreateAt,
			InviterId:        app.InviterId,
			OperatorId:       app.OperatorId,
			Status:           group.ApplicationStatus(app.Status),
			UserNickname:     user.Nickname,
			UserAvatar:       user.Avatar,
			GroupName:        grouper.Name,
			GroupAvatar:      grouper.Avatar,
			InviteNickname:   "",
			InviteAvatar:     "",
			OperatorNickname: "",
			OperatorAvatar:   "",
		}

		if inviter != nil {
			pbApp.InviteNickname = inviter.Nickname
			pbApp.InviteAvatar = inviter.Avatar
		}
		if operator != nil {
			pbApp.OperatorNickname = operator.Nickname
			pbApp.OperatorAvatar = operator.Avatar
		}

		pbApplications = append(pbApplications, pbApp)
	}

	return &group.GetJoinGroupApplicationsResponse{
		Applications: pbApplications,
	}, nil
}
