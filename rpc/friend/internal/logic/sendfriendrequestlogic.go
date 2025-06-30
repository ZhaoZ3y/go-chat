package logic

import (
	"IM/pkg/model"
	"context"
	"time"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFriendRequestLogic {
	return &SendFriendRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送好友申请
func (l *SendFriendRequestLogic) SendFriendRequest(in *friend.SendFriendRequestRequest) (*friend.SendFriendRequestResponse, error) {
	// 验证参数
	if in.FromUserId == 0 || in.ToUserId == 0 {
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "用户ID不能为空",
		}, nil
	}

	if in.FromUserId == in.ToUserId {
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "不能添加自己为好友",
		}, nil
	}

	// 使用事务确保数据一致性
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已经是好友
	var friendCount int64
	err := tx.Model(&model.Friends{}).
		Where("user_id = ? AND friend_id = ? AND status = 1", in.FromUserId, in.ToUserId).
		Count(&friendCount).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("查询好友关系失败: %v", err)
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "系统错误，请稍后重试",
		}, nil
	}

	if friendCount > 0 {
		tx.Rollback()
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "已经是好友关系",
		}, nil
	}

	// 检查是否已有待处理的申请
	var requestCount int64
	err = tx.Model(&model.FriendRequests{}).
		Where("from_user_id = ? AND to_user_id = ? AND status = 1", in.FromUserId, in.ToUserId).
		Count(&requestCount).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("查询好友申请失败: %v", err)
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "系统错误，请稍后重试",
		}, nil
	}

	if requestCount > 0 {
		tx.Rollback()
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "已发送好友申请，请等待对方处理",
		}, nil
	}

	// 创建好友申请
	request := &model.FriendRequests{
		FromUserId: in.FromUserId,
		ToUserId:   in.ToUserId,
		Message:    in.Message,
		Status:     1, // 待处理
		CreateAt:   time.Now().Unix(),
		UpdateAt:   time.Now().Unix(),
	}

	if err := tx.Create(request).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("创建好友申请失败: %v", err)
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "发送好友申请失败",
		}, nil
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "发送好友申请失败",
		}, nil
	}

	return &friend.SendFriendRequestResponse{
		Success:   true,
		Message:   "好友申请发送成功",
		RequestId: request.Id, // 返回申请ID，便于后续操作
	}, nil
}
