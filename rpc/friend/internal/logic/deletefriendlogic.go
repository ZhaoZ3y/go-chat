package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除好友
func (l *DeleteFriendLogic) DeleteFriend(in *friend.DeleteFriendRequest) (*friend.DeleteFriendResponse, error) {
	// 验证参数
	if in.UserId == 0 || in.FriendId == 0 {
		return &friend.DeleteFriendResponse{
			Success: false,
			Message: "参数错误",
		}, nil
	}

	// 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除双向好友关系
	if err := tx.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).Delete(&model.Friends{}).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除好友关系失败: %v", err)
		return &friend.DeleteFriendResponse{
			Success: false,
			Message: "删除好友失败",
		}, nil
	}

	if err := tx.Where("user_id = ? AND friend_id = ?", in.FriendId, in.UserId).Delete(&model.Friends{}).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除好友关系失败: %v", err)
		return &friend.DeleteFriendResponse{
			Success: false,
			Message: "删除好友失败",
		}, nil
	}

	if err := tx.Commit().Error; err != nil {
		return &friend.DeleteFriendResponse{
			Success: false,
			Message: "删除好友失败",
		}, nil
	}

	return &friend.DeleteFriendResponse{
		Success: true,
		Message: "删除好友成功",
	}, nil
}
