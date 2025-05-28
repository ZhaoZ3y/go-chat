package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFriendRequestLogic {
	return &HandleFriendRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 处理好友申请
func (l *HandleFriendRequestLogic) HandleFriendRequest(in *friend.HandleFriendRequestRequest) (*friend.HandleFriendRequestResponse, error) {
	// 验证参数
	if in.RequestId == 0 || in.UserId == 0 {
		return &friend.HandleFriendRequestResponse{
			Success: false,
			Message: "参数错误",
		}, nil
	}

	if in.Action != 2 && in.Action != 3 {
		return &friend.HandleFriendRequestResponse{
			Success: false,
			Message: "操作类型错误",
		}, nil
	}

	// 查找好友申请
	var request model.FriendRequests
	err := l.svcCtx.DB.Where("id = ? AND to_user_id = ? AND status = 1", in.RequestId, in.UserId).
		First(&request).Error
	if err != nil {
		return &friend.HandleFriendRequestResponse{
			Success: false,
			Message: "好友申请不存在或已处理",
		}, nil
	}

	// 开启事务处理
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新申请状态
	if err := tx.Model(&request).Update("status", in.Action).Error; err != nil {
		tx.Rollback()
		return &friend.HandleFriendRequestResponse{
			Success: false,
			Message: "处理申请失败",
		}, nil
	}

	// 如果同意申请，创建双向好友关系
	if in.Action == 2 {
		// 检查是否已经是好友关系
		var existFriend model.Friends
		err := tx.Where("user_id = ? AND friend_id = ?", request.FromUserId, request.ToUserId).
			First(&existFriend).Error
		if err == nil {
			// 已存在好友关系，只需更新状态
			if err := tx.Model(&existFriend).Update("status", 1).Error; err != nil {
				tx.Rollback()
				return &friend.HandleFriendRequestResponse{
					Success: false,
					Message: "建立好友关系失败",
				}, nil
			}
		} else {
			// 创建申请者到接受者的好友关系
			friend1 := &model.Friends{
				UserId:   request.FromUserId,
				FriendId: request.ToUserId,
				Status:   1,
			}
			if err := tx.Create(friend1).Error; err != nil {
				tx.Rollback()
				return &friend.HandleFriendRequestResponse{
					Success: false,
					Message: "建立好友关系失败",
				}, nil
			}
		}

		// 检查反向好友关系
		var existFriend2 model.Friends
		err = tx.Where("user_id = ? AND friend_id = ?", request.ToUserId, request.FromUserId).
			First(&existFriend2).Error
		if err == nil {
			// 已存在好友关系，只需更新状态
			if err := tx.Model(&existFriend2).Update("status", 1).Error; err != nil {
				tx.Rollback()
				return &friend.HandleFriendRequestResponse{
					Success: false,
					Message: "建立好友关系失败",
				}, nil
			}
		} else {
			// 创建接受者到申请者的好友关系
			friend2 := &model.Friends{
				UserId:   request.ToUserId,
				FriendId: request.FromUserId,
				Status:   1,
			}
			if err := tx.Create(friend2).Error; err != nil {
				tx.Rollback()
				return &friend.HandleFriendRequestResponse{
					Success: false,
					Message: "建立好友关系失败",
				}, nil
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return &friend.HandleFriendRequestResponse{
			Success: false,
			Message: "处理申请失败",
		}, nil
	}

	message := "已拒绝好友申请"
	if in.Action == 2 {
		message = "已同意好友申请"
	}

	return &friend.HandleFriendRequestResponse{
		Success: true,
		Message: message,
	}, nil
}
