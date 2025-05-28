package logic

import (
	"IM/pkg/model"
	"context"
	"gorm.io/gorm"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickFromGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickFromGroupLogic {
	return &KickFromGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 踢出群组
func (l *KickFromGroupLogic) KickFromGroup(in *group.KickFromGroupRequest) (*group.KickFromGroupResponse, error) {
	// 检查操作者权限
	var operatorMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role IN (1,2)",
		in.GroupId, in.OperatorId).First(&operatorMember).Error; err != nil {
		return &group.KickFromGroupResponse{Success: false, Message: "无权限操作"}, nil
	}

	// 检查被踢用户
	var targetMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&targetMember).Error; err != nil {
		return &group.KickFromGroupResponse{Success: false, Message: "用户不在群组中"}, nil
	}

	// 不能踢出群主
	if targetMember.Role == 1 {
		return &group.KickFromGroupResponse{Success: false, Message: "不能踢出群主"}, nil
	}

	// 管理员不能踢出其他管理员
	if operatorMember.Role == 2 && targetMember.Role == 2 {
		return &group.KickFromGroupResponse{Success: false, Message: "管理员不能踢出其他管理员"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	if err := tx.Delete(&targetMember).Error; err != nil {
		tx.Rollback()
		return &group.KickFromGroupResponse{Success: false, Message: "踢出失败"}, nil
	}

	tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
		Update("member_count", gorm.Expr("member_count - 1"))

	tx.Commit()
	return &group.KickFromGroupResponse{Success: true, Message: "踢出成功"}, nil
}
