package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendRemarkLogic {
	return &UpdateFriendRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新好友备注
func (l *UpdateFriendRemarkLogic) UpdateFriendRemark(in *friend.UpdateFriendRemarkRequest) (*friend.UpdateFriendRemarkResponse, error) {
	// 验证参数
	if in.UserId == 0 || in.FriendId == 0 {
		return &friend.UpdateFriendRemarkResponse{
			Success: false,
			Message: "参数错误",
		}, nil
	}

	// 查找好友关系
	var friendRelation model.Friends
	err := l.svcCtx.DB.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).First(&friendRelation).Error
	if err != nil {
		return &friend.UpdateFriendRemarkResponse{
			Success: false,
			Message: "好友关系不存在",
		}, nil
	}

	// 更新备注
	if err := l.svcCtx.DB.Model(&friendRelation).Update("remark", in.Remark).Error; err != nil {
		l.Logger.Errorf("更新好友备注失败: %v", err)
		return &friend.UpdateFriendRemarkResponse{
			Success: false,
			Message: "更新好友备注失败",
		}, nil
	}

	return &friend.UpdateFriendRemarkResponse{
		Success: true,
		Message: "更新好友备注成功",
	}, nil
}
