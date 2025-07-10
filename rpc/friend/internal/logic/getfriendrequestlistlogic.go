package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendRequestListLogic {
	return &GetFriendRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取好友申请列表
func (l *GetFriendRequestListLogic) GetFriendRequestList(in *friend.GetFriendRequestListRequest) (*friend.GetFriendRequestListResponse, error) {
	// 验证参数
	if in.UserId == 0 {
		return &friend.GetFriendRequestListResponse{}, nil
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&model.FriendRequests{}).Where("to_user_id = ?", in.UserId)

	if in.Status > 0 {
		query = query.Where("status = ?", in.Status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("获取好友申请总数失败: %v", err)
		return &friend.GetFriendRequestListResponse{}, nil
	}

	// 获取列表数据
	var requests []model.FriendRequests
	if err := query.Order("create_at DESC").Find(&requests).Error; err != nil {
		l.Logger.Errorf("获取好友申请列表失败: %v", err)
		return &friend.GetFriendRequestListResponse{}, nil
	}

	// 转换为响应格式（批量优化）
	var result []*friend.FriendRequest

	// 1. 收集所有相关用户ID
	userIDSet := make(map[int64]struct{})
	for _, req := range requests {
		userIDSet[req.FromUserId] = struct{}{}
		userIDSet[req.ToUserId] = struct{}{}
	}

	var userIds []int64
	for id := range userIDSet {
		userIds = append(userIds, id)
	}

	// 2. 批量查询用户信息
	var users []model.User
	if err := l.svcCtx.DB.Where("id IN ?", userIds).Find(&users).Error; err != nil {
		l.Logger.Errorf("批量查询用户信息失败: %v", err)
		return &friend.GetFriendRequestListResponse{}, nil
	}

	// 3. 构建 userId => User 映射
	userMap := make(map[int64]*model.User)
	for _, u := range users {
		userCopy := u // 避免引用错误
		userMap[u.Id] = &userCopy
	}

	// 4. 构建响应
	for _, req := range requests {
		fromUser, ok1 := userMap[req.FromUserId]
		toUser, ok2 := userMap[req.ToUserId]
		if !ok1 || !ok2 {
			l.Logger.Errorf("用户信息缺失，from: %d, to: %d", req.FromUserId, req.ToUserId)
			continue
		}

		result = append(result, &friend.FriendRequest{
			Id:           req.Id,
			FromUserId:   req.FromUserId,
			ToUserId:     req.ToUserId,
			Message:      req.Message,
			Status:       int32(req.Status),
			CreateAt:     req.CreateAt,
			UpdateAt:     req.UpdateAt,
			FromNickname: fromUser.Nickname,
			FromAvatar:   fromUser.Avatar,
			ToNickname:   toUser.Nickname,
			ToAvatar:     toUser.Avatar,
		})
	}

	return &friend.GetFriendRequestListResponse{
		Requests: result,
		Total:    total,
	}, nil
}
