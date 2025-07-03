package logic

import (
	"IM/pkg/model"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupInfoLogic {
	return &UpdateGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新群组信息
func (l *UpdateGroupInfoLogic) UpdateGroupInfo(in *group.UpdateGroupInfoRequest) (*group.UpdateGroupInfoResponse, error) {
	if in.GroupId == 0 || in.OperatorId == 0 {
		return &group.UpdateGroupInfoResponse{Success: false, Message: "参数错误"}, nil
	}

	var operatorMember model.GroupMembers
	err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.OperatorId).First(&operatorMember).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.UpdateGroupInfoResponse{Success: false, Message: "您不是该群成员，无权操作"}, nil
		}
		l.Logger.Errorf("UpdateGroupInfo: find operator member failed: %v", err)
		return &group.UpdateGroupInfoResponse{Success: false, Message: "查询权限信息失败"}, nil
	}

	if operatorMember.Role != int64(group.MemberRole_ROLE_OWNER) && operatorMember.Role != int64(group.MemberRole_ROLE_ADMIN) {
		return &group.UpdateGroupInfoResponse{Success: false, Message: "只有群主或管理员才能修改群信息"}, nil
	}

	updateData := make(map[string]interface{})
	if in.Name != "" {
		updateData["name"] = in.Name
	}
	if in.Description != "" {
		updateData["description"] = in.Description
	}
	if in.Avatar != "" {
		updateData["avatar"] = in.Avatar
	}

	if len(updateData) == 0 {
		return &group.UpdateGroupInfoResponse{Success: true, Message: "没有需要更新的信息"}, nil
	}

	dbResult := l.svcCtx.DB.Model(&model.Groups{}).Where("id = ?", in.GroupId).Updates(updateData)
	if dbResult.Error != nil {
		l.Logger.Errorf("UpdateGroupInfo: update group info failed: %v", dbResult.Error)
		return &group.UpdateGroupInfoResponse{Success: false, Message: "更新群组信息失败"}, nil
	}

	return &group.UpdateGroupInfoResponse{
		Success: true,
		Message: "群组信息更新成功",
	}, nil
}
