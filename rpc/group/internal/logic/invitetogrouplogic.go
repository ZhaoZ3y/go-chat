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

type InviteToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInviteToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteToGroupLogic {
	return &InviteToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 邀请加入群组
func (l *InviteToGroupLogic) InviteToGroup(in *group.InviteToGroupRequest) (*group.InviteToGroupResponse, error) {
	// 检查邀请者权限
	var inviterMember model.GroupMembers
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND role IN (1,2)",
		in.GroupId, in.InviterId).First(&inviterMember).Error; err != nil {
		return &group.InviteToGroupResponse{Success: false, Message: "无权限邀请"}, nil
	}

	var failedUserIds []int64
	tx := l.svcCtx.DB.Begin()

	for _, userId := range in.UserIds {
		// 检查是否已经是成员
		var existMember model.GroupMembers
		if err := tx.Where("group_id = ? AND user_id = ?", in.GroupId, userId).First(&existMember).Error; err == nil {
			failedUserIds = append(failedUserIds, userId)
			continue
		}

		// 添加新成员
		member := &model.GroupMembers{
			GroupId:  in.GroupId,
			UserId:   userId,
			Role:     3,
			Status:   1,
			JoinTime: time.Now().Unix(),
		}

		if err := tx.Create(member).Error; err != nil {
			failedUserIds = append(failedUserIds, userId)
			continue
		}

		// 更新群成员数量
		tx.Model(&model.Groups{}).Where("id = ?", in.GroupId).
			Update("member_count", gorm.Expr("member_count + 1"))
	}

	tx.Commit()
	return &group.InviteToGroupResponse{
		Success:       true,
		Message:       "邀请完成",
		FailedUserIds: failedUserIds,
	}, nil
}
