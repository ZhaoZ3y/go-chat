package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockFriendLogic {
	return &BlockFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 拉黑好友
func (l *BlockFriendLogic) BlockFriend(in *friend.BlockFriendRequest) (*friend.BlockFriendResponse, error) {
	// 验证参数
	if in.UserId == 0 || in.FriendId == 0 {
		return &friend.BlockFriendResponse{
			Success: false,
			Message: "参数错误",
		}, nil
	}

	// 查找好友关系
	var friendRelation model.Friends
	err := l.svcCtx.DB.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).First(&friendRelation).Error
	if err != nil {
		return &friend.BlockFriendResponse{
			Success: false,
			Message: "好友关系不存在",
		}, nil
	}

	// 更新状态为拉黑
	if err := l.svcCtx.DB.Model(&friendRelation).Update("status", 2).Error; err != nil {
		l.Logger.Errorf("拉黑好友失败: %v", err)
		return &friend.BlockFriendResponse{
			Success: false,
			Message: "拉黑好友失败",
		}, nil
	}

	return &friend.BlockFriendResponse{
		Success: true,
		Message: "拉黑好友成功",
	}, nil
}
