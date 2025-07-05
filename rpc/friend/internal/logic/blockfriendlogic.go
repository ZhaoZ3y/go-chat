package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/const"
	"context"
	"gorm.io/gorm"

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
	if in.UserId == 0 || in.FriendId == 0 {
		return &friend.BlockFriendResponse{Success: false, Message: "用户ID和好友ID不能为空"}, nil
	}
	if in.UserId == in.FriendId {
		return &friend.BlockFriendResponse{Success: false, Message: "不能操作自己"}, nil
	}

	// 1. 查找我与好友的关系
	var myFriendRelation model.Friends
	err := l.svcCtx.DB.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).First(&myFriendRelation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &friend.BlockFriendResponse{Success: false, Message: "好友关系不存在"}, nil
		}
		l.Logger.Errorf("查询好友关系失败: %v", err)
		return &friend.BlockFriendResponse{Success: false, Message: "操作失败，请稍后重试"}, nil
	}

	// 2. 根据当前状态决定是拉黑还是取消拉黑
	var newMyStatus, newFriendStatus int
	var responseMessage string

	if myFriendRelation.Status == _const.FriendStatusBlocked {
		// 当前已拉黑，执行“取消拉黑”操作
		newMyStatus = _const.FriendStatusNormal
		newFriendStatus = _const.FriendStatusNormal
		responseMessage = "已取消拉黑"
	} else if myFriendRelation.Status == _const.FriendStatusNormal {
		// 当前是正常好友，执行“拉黑”操作
		newMyStatus = _const.FriendStatusBlocked
		newFriendStatus = _const.FriendStatusBeBlocked
		responseMessage = "拉黑好友成功"
	} else {
		// 比如对方已经拉黑了你 (FriendStatusBeBlocked)，你不应该能直接操作
		return &friend.BlockFriendResponse{Success: false, Message: "当前状态无法进行此操作"}, nil
	}

	// 3. 开启事务执行更新
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新我的关系
	if err := tx.Model(&model.Friends{}).Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).Update("status", newMyStatus).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新我的好友状态失败: %v", err)
		return &friend.BlockFriendResponse{Success: false, Message: "操作失败，请稍后重试"}, nil
	}

	// 更新对方的关系
	if err := tx.Model(&model.Friends{}).Where("user_id = ? AND friend_id = ?", in.FriendId, in.UserId).Update("status", newFriendStatus).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新对方的好友状态失败: %v", err)
		return &friend.BlockFriendResponse{Success: false, Message: "操作失败，请稍后重试"}, nil
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("拉黑/取消拉黑事务提交失败: %v", err)
		return &friend.BlockFriendResponse{Success: false, Message: "操作失败，请稍后重试"}, nil
	}

	return &friend.BlockFriendResponse{Success: true, Message: responseMessage}, nil
}
