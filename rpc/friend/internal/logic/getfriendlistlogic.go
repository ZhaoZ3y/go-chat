package logic

import (
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取好友列表
func (l *GetFriendListLogic) GetFriendList(in *friend.GetFriendListRequest) (*friend.GetFriendListResponse, error) {
	// 验证参数
	if in.UserId == 0 {
		return &friend.GetFriendListResponse{}, nil
	}

	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}

	// 构建查询条件，只查询正常状态的好友
	query := l.svcCtx.DB.Model(&model.Friends{}).Where("user_id = ? AND status = 1", in.UserId)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("获取好友总数失败: %v", err)
		return &friend.GetFriendListResponse{}, nil
	}

	// 获取列表数据
	var friends []model.Friends
	offset := (in.Page - 1) * in.PageSize
	if err := query.Order("create_at DESC").Offset(int(offset)).Limit(int(in.PageSize)).Find(&friends).Error; err != nil {
		l.Logger.Errorf("获取好友列表失败: %v", err)
		return &friend.GetFriendListResponse{}, nil
	}

	// 转换为响应格式
	var result []*friend.Friend
	for _, f := range friends {
		result = append(result, &friend.Friend{
			Id:       f.Id,
			UserId:   f.UserId,
			FriendId: f.FriendId,
			Remark:   f.Remark,
			Status:   int32(f.Status),
			CreateAt: f.CreateAt,
			UpdateAt: f.UpdateAt,
		})
	}

	return &friend.GetFriendListResponse{
		Friends: result,
		Total:   total,
	}, nil
}
