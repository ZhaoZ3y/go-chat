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
	var applications []*model.JoinGroupApplications
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("to_group_id = ?", in.GroupId).
		Order("create_at DESC").
		Find(&applications).Error

	if err != nil {
		l.Logger.Errorf("failed to get join group applications for group %d: %v", in.GroupId, err)
		return nil, err
	}

	var pbApplications []*group.JoinGroupApplication
	for _, application := range applications {
		pbApplication := &group.JoinGroupApplication{
			Id:        application.Id,
			UserId:    application.FromUserId,
			GroupId:   application.ToGroupId,
			Reason:    application.Reason,
			InviterId: application.InviterId,
			Status:    group.ApplicationStatus(application.Status),
			ApplyTime: application.CreateAt,
		}
		pbApplications = append(pbApplications, pbApplication)
	}

	return &group.GetJoinGroupApplicationsResponse{
		Applications: pbApplications,
	}, nil
}
