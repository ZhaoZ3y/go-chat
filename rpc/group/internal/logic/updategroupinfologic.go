package logic

import (
	"IM/pkg/model"
	"context"

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
	// 检查操作者权限
	var member model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role IN (1,2)",
		in.GroupId, in.OperatorId).First(&member).Error; err != nil {
		return &group.UpdateGroupInfoResponse{Success: false, Message: "无权限操作"}, nil
	}

	updates := make(map[string]interface{})
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Description != "" {
		updates["description"] = in.Description
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}

	if err := l.svcCtx.DB.Model(&model.Groups{}).Where("id = ?", in.GroupId).Updates(updates).Error; err != nil {
		return &group.UpdateGroupInfoResponse{Success: false, Message: "更新失败"}, nil
	}

	return &group.UpdateGroupInfoResponse{Success: true, Message: "更新成功"}, nil
}
