package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetMemberRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMemberRoleLogic {
	return &SetMemberRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设置群组成员角色
func (l *SetMemberRoleLogic) SetMemberRole(in *group.SetMemberRoleRequest) (*group.SetMemberRoleResponse, error) {
	// 检查操作者权限（只有群主可以设置角色）
	var operatorMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role = 1",
		in.GroupId, in.OperatorId).First(&operatorMember).Error; err != nil {
		return &group.SetMemberRoleResponse{Success: false, Message: "只有群主可以设置角色"}, nil
	}

	// 更新目标用户角色
	if err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
		Update("role", in.Role).Error; err != nil {
		return &group.SetMemberRoleResponse{Success: false, Message: "设置角色失败"}, nil
	}

	return &group.SetMemberRoleResponse{Success: true, Message: "设置角色成功"}, nil
}
