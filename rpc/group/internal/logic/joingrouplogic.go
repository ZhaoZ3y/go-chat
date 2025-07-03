package logic

import (
	"IM/pkg/model"
	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"

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
	if in.GroupId == 0 || in.UserId == 0 {
		return &group.JoinGroupResponse{Success: false, Message: "参数错误：群组ID和用户ID不能为空"}, nil
	}

	var targetGroup model.Groups
	err := l.svcCtx.DB.Where("id = ? AND status = ?", in.GroupId, group.GroupStatus_GROUP_STATUS_NORMAL).First(&targetGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.JoinGroupResponse{Success: false, Message: "群组不存在或已解散"}, nil
		}
		l.Logger.Errorf("JoinGroup: find group failed, GroupID: %d, Error: %v", in.GroupId, err)
		return &group.JoinGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 检查用户是否已经是群成员
	var memberCount int64
	l.svcCtx.DB.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Count(&memberCount)
	if memberCount > 0 {
		return &group.JoinGroupResponse{Success: false, Message: "您已是该群组成员"}, nil
	}

	// 检查是否已有待处理的申请
	var applicationCount int64
	l.svcCtx.DB.Model(&model.JoinGroupApplications{}).Where("to_group_id = ? AND from_user_id = ? AND status = ?", in.GroupId, in.UserId, group.ApplicationStatus_PENDING).Count(&applicationCount)
	if applicationCount > 0 {
		return &group.JoinGroupResponse{Success: false, Message: "您已提交过申请，请耐心等待管理员审核"}, nil
	}

	application := &model.JoinGroupApplications{
		FromUserId: in.UserId,
		ToGroupId:  in.GroupId,
		Reason:     in.Reason,
		InviterId:  0, // 0 表示主动申请
		Status:     int8(group.ApplicationStatus_PENDING),
	}

	if err := l.svcCtx.DB.Create(application).Error; err != nil {
		l.Logger.Errorf("JoinGroup: create application failed: %v", err)
		return &group.JoinGroupResponse{Success: false, Message: "提交申请失败，请稍后再试"}, nil
	}

	//TODO：后续添加消息队列异步发送站外通知，类型为 NOTIFY_MEMBER_APPLY_JOIN

	return &group.JoinGroupResponse{
		Success: true,
		Message: "入群申请已提交，请等待管理员审核",
	}, nil
}
