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

type DismissGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDismissGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DismissGroupLogic {
	return &DismissGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解散群组
func (l *DismissGroupLogic) DismissGroup(in *group.DismissGroupRequest) (*group.DismissGroupResponse, error) {
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		logx.Errorf("DismissGroup: begin transaction failed, error: %v", tx.Error)
		return nil, tx.Error
	}
	// 确保在函数退出时，如果事务未提交，则回滚
	defer tx.Rollback()

	var groupModel model.Groups
	if err := tx.Where("id = ?", in.GroupId).First(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &group.DismissGroupResponse{Success: false, Message: "群组不存在"}, nil
		}
		logx.Errorf("DismissGroup: find group failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "查询群组信息失败"}, nil
	}

	// 验证操作者是否为群主
	if groupModel.OwnerId != in.OwnerId {
		return &group.DismissGroupResponse{Success: false, Message: "无权解散群组，操作者非群主"}, nil
	}

	// 删除所有群组成员
	if err := tx.Where("group_id = ?", in.GroupId).Delete(&model.GroupMembers{}).Error; err != nil {
		logx.Errorf("DismissGroup: delete group members failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "删除群组成员失败"}, nil
	}

	// 删除所有相关的入群申请
	if err := tx.Where("to_group_id = ?", in.GroupId).Delete(&model.JoinGroupApplications{}).Error; err != nil {
		logx.Errorf("DismissGroup: delete join applications failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "清理入群申请失败"}, nil
	}

	// 删除群组本身
	if err := tx.Delete(&groupModel).Error; err != nil {
		logx.Errorf("DismissGroup: delete group record failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "删除群组记录失败"}, nil
	}

	//
	// TODO: 在这里可以添加逻辑，将解散通知 (NOTIFY_GROUP_DISMISSED) 推送给所有原成员。
	// 最好在事务提交后异步处理，以确保操作最终成功。
	// 可以先查询出所有成员ID，然后在事务提交后，通过消息队列分发通知任务。
	//

	if err := tx.Commit().Error; err != nil {
		logx.Errorf("DismissGroup: commit transaction failed, group_id: %d, error: %v", in.GroupId, err)
		return &group.DismissGroupResponse{Success: false, Message: "解散群组事务提交失败"}, nil
	}

	return &group.DismissGroupResponse{
		Success: true,
		Message: "群组已成功解散",
	}, nil
}
