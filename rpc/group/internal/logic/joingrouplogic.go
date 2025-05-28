package logic

import (
	"IM/pkg/model"
	"context"
	"gorm.io/gorm"
	"time"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 加入群组
func (l *JoinGroupLogic) JoinGroup(in *group.JoinGroupRequest) (*group.JoinGroupResponse, error) {
	// 检查群组是否存在
	var groupInfo model.Groups
	if err := l.svcCtx.DB.Where("id = ? AND status = 1", in.GroupId).First(&groupInfo).Error; err != nil {
		return &group.JoinGroupResponse{Success: false, Message: "群组不存在"}, nil
	}

	// 检查是否已经是成员
	var existMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&existMember).Error; err == nil {
		return &group.JoinGroupResponse{Success: false, Message: "已经是群成员"}, nil
	}

	// 检查群成员数量限制
	if groupInfo.MemberCount >= groupInfo.MaxMemberCount {
		return &group.JoinGroupResponse{Success: false, Message: "群组成员已满"}, nil
	}

	// 添加新成员
	tx := l.svcCtx.DB.Begin()
	member := &model.GroupMembers{
		GroupId:  in.GroupId,
		UserId:   in.UserId,
		Role:     3, // 普通成员
		Status:   1,
		JoinTime: time.Now().Unix(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return &group.JoinGroupResponse{Success: false, Message: "加入群组失败"}, nil
	}

	// 更新群成员数量
	tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
		Update("member_count", gorm.Expr("member_count + 1"))

	tx.Commit()
	return &group.JoinGroupResponse{Success: true, Message: "加入群组成功"}, nil
}
