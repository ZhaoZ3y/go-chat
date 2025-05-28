package logic

import (
	"IM/pkg/model"
	"context"

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

	// 检查是否已经是好友
	var existFriend model.Friends
	err := l.svcCtx.DB.Where("user_id = ? AND friend_id = ? AND status = 1", in.FromUserId, in.ToUserId).
		First(&existFriend).Error
	if err == nil {
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "已经是好友关系",
		}, nil
	}

	// 检查是否已有待处理的申请
	var existRequest model.FriendRequests
	err = l.svcCtx.DB.Where("from_user_id = ? AND to_user_id = ? AND status = 1", in.FromUserId, in.ToUserId).
		First(&existRequest).Error
	if err == nil {
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
	}

	if err := l.svcCtx.DB.Create(request).Error; err != nil {
		l.Logger.Errorf("创建好友申请失败: %v", err)
		return &friend.SendFriendRequestResponse{
			Success: false,
			Message: "发送好友申请失败",
		}, nil
	}

	return &friend.SendFriendRequestResponse{
		Success: true,
		Message: "好友申请发送成功",
	}, nil
}
