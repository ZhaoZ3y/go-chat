package logic

import (
	"IM/pkg/model"
	"context"
	"errors"
	"gorm.io/gorm"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群组信息
func (l *GetGroupInfoLogic) GetGroupInfo(in *group.GetGroupInfoRequest) (*group.GetGroupInfoResponse, error) {
	var memberModel model.GroupMembers
	err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&memberModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.GetGroupInfoResponse{}, nil
		}
		logx.Errorf("GetGroupInfo: find group member failed, group_id: %d, user_id: %d, error: %v", in.GroupId, in.UserId, err)
		return nil, err
	}

	var groupModel model.Groups
	if err := l.svcCtx.DB.Where("id = ?", in.GroupId).First(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("GetGroupInfo: data inconsistency, group not found but member exists, group_id: %d", in.GroupId)
			return &group.GetGroupInfoResponse{}, nil
		}
		logx.Errorf("GetGroupInfo: find group failed, group_id: %d, error: %v", in.GroupId, err)
		return nil, err
	}

	// 转换群组信息
	groupInfoPb := &group.Group{
		Id:             groupModel.Id,
		Name:           groupModel.Name,
		Description:    groupModel.Description,
		Avatar:         groupModel.Avatar,
		OwnerId:        groupModel.OwnerId,
		MemberCount:    int32(groupModel.MemberCount),
		MaxMemberCount: int32(groupModel.MaxMemberCount),
		Status:         group.GroupStatus(groupModel.Status),
		CreateAt:       groupModel.CreateAt,
		UpdateAt:       groupModel.UpdateAt,
	}

	return &group.GetGroupInfoResponse{
		GroupInfo: groupInfoPb,
	}, nil
}
