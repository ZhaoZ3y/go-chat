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

type UpdateGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupMemberInfoLogic {
	return &UpdateGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateGroupMemberInfo 修改群组成员信息 (例如：修改自己的群昵称)
func (l *UpdateGroupMemberInfoLogic) UpdateGroupMemberInfo(in *group.UpdateGroupMemberInfoRequest) (*group.UpdateGroupMemberInfoResponse, error) {
	if in.GroupId == 0 || in.UserId == 0 {
		return &group.UpdateGroupMemberInfoResponse{Success: false, Message: "参数错误：群组ID和用户ID不能为空"}, nil
	}
	if in.Nickname == "" {
		return &group.UpdateGroupMemberInfoResponse{Success: false, Message: "昵称不能为空"}, nil
	}

	var member model.GroupMembers
	err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.UpdateGroupMemberInfoResponse{Success: false, Message: "您不是该群组成员，无法修改信息"}, nil
		}
		l.Logger.Errorf("UpdateGroupMemberInfo: find member failed: %v", err)
		return &group.UpdateGroupMemberInfoResponse{Success: false, Message: "查询成员信息失败"}, nil
	}

	// 如果昵称没有变化，直接返回成功
	if member.Nickname == in.Nickname {
		return &group.UpdateGroupMemberInfoResponse{Success: true, Message: "昵称未发生变化"}, nil
	}

	err = l.svcCtx.DB.Model(&model.GroupMembers{}).
		Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).
		Update("nickname", in.Nickname).Error

	if err != nil {
		l.Logger.Errorf("UpdateGroupMemberInfo: update member nickname failed: %v", err)
		return &group.UpdateGroupMemberInfoResponse{Success: false, Message: "更新昵称失败"}, nil
	}

	return &group.UpdateGroupMemberInfoResponse{
		Success: true,
		Message: "群昵称更新成功",
	}, nil
}
