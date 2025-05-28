package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MuteMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMuteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteMemberLogic {
	return &MuteMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁言群组成员
func (l *MuteMemberLogic) MuteMember(in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	// 检查操作者权限
	var operatorMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role IN (1,2)",
		in.GroupId, in.OperatorId).First(&operatorMember).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "无权限操作"}, nil
	}

	// 检查目标用户
	var targetMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&targetMember).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "用户不在群组中"}, nil
	}

	// 不能禁言群主
	if targetMember.Role == 1 {
		return &group.MuteMemberResponse{Success: false, Message: "不能禁言群主"}, nil
	}

	// 管理员不能禁言其他管理员
	if operatorMember.Role == 2 && targetMember.Role == 2 {
		return &group.MuteMemberResponse{Success: false, Message: "管理员不能禁言其他管理员"}, nil
	}

	// 设置禁言状态
	status := int8(2) // 禁言
	if in.Duration == 0 {
		status = 1 // 取消禁言
	}

	if err := l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
		Update("status", status).Error; err != nil {
		return &group.MuteMemberResponse{Success: false, Message: "操作失败"}, nil
	}

	message := "禁言成功"
	if in.Duration == 0 {
		message = "取消禁言成功"
	}

	return &group.MuteMemberResponse{Success: true, Message: message}, nil
}
