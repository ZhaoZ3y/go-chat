package logic

import (
	"IM/pkg/model"
	"context"
	"gorm.io/gorm"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupLogic {
	return &LeaveGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 退出群组
func (l *LeaveGroupLogic) LeaveGroup(in *group.LeaveGroupRequest) (*group.LeaveGroupResponse, error) {
	// 检查是否为群主
	var member model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&member).Error; err != nil {
		return &group.LeaveGroupResponse{Success: false, Message: "不是群成员"}, nil
	}

	if member.Role == 1 {
		return &group.LeaveGroupResponse{Success: false, Message: "群主不能退出群组，请先转让群组"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	// 删除成员记录
	if err := tx.Delete(&member).Error; err != nil {
		tx.Rollback()
		return &group.LeaveGroupResponse{Success: false, Message: "退出群组失败"}, nil
	}

	// 更新群成员数量
	tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
		Update("member_count", gorm.Expr("member_count - 1"))

	tx.Commit()
	return &group.LeaveGroupResponse{Success: true, Message: "退出群组成功"}, nil
}
