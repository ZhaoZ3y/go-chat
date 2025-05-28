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

	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
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
	offset := (in.Page - 1) * in.PageSize
	if err := query.Order("create_at DESC").Offset(int(offset)).Limit(int(in.PageSize)).Find(&requests).Error; err != nil {
		l.Logger.Errorf("获取好友申请列表失败: %v", err)
		return &friend.GetFriendRequestListResponse{}, nil
	}

	// 转换为响应格式
	var result []*friend.FriendRequest
	for _, req := range requests {
		result = append(result, &friend.FriendRequest{
			Id:         req.Id,
			FromUserId: req.FromUserId,
			ToUserId:   req.ToUserId,
			Message:    req.Message,
			Status:     int32(req.Status),
			CreateAt:   req.CreateAt,
			UpdateAt:   req.UpdateAt,
		})
	}

	return &friend.GetFriendRequestListResponse{
		Requests: result,
		Total:    total,
	}, nil
}
