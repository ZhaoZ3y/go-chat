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

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建群组
func (l *CreateGroupLogic) CreateGroup(in *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建群组
	groupModel := &model.Groups{
		Name:           in.Name,
		Description:    in.Description,
		Avatar:         in.Avatar,
		OwnerId:        in.OwnerId,
		MemberCount:    1,
		MaxMemberCount: 500,
		Status:         1,
	}

	if err := tx.Create(groupModel).Error; err != nil {
		tx.Rollback()
		return &group.CreateGroupResponse{Success: false, Message: "创建群组失败"}, nil
	}

	// 添加群主为成员
	ownerMember := &model.GroupMembers{
		GroupId:  groupModel.Id,
		UserId:   in.OwnerId,
		Role:     1, // 群主
		Status:   1,
		JoinTime: time.Now().Unix(),
	}

	if err := tx.Create(ownerMember).Error; err != nil {
		tx.Rollback()
		return &group.CreateGroupResponse{Success: false, Message: "添加群主失败"}, nil
	}

	// 添加其他成员
	for _, userId := range in.MemberIds {
		if userId == in.OwnerId {
			continue
		}

		member := &model.GroupMembers{
			GroupId:  groupModel.Id,
			UserId:   userId,
			Role:     3, // 普通成员
			Status:   1,
			JoinTime: time.Now().Unix(),
		}

		if err := tx.Create(member).Error; err != nil {
			continue // 跳过失败的用户
		}

		// 更新群成员数量
		tx.Model(&model.Groups{}).Where("id = ?", groupModel.Id).
			Update("member_count", gorm.Expr("member_count + 1"))
	}

	tx.Commit()
	return &group.CreateGroupResponse{
		GroupId: groupModel.Id,
		Success: true,
		Message: "创建群组成功",
	}, nil
}
