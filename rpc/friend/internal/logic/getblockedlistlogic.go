package logic

import (
	_const "IM/pkg/const"
	"IM/pkg/model"
	"context"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlockedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlockedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockedListLogic {
	return &GetBlockedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取拉黑列表
func (l *GetBlockedListLogic) GetBlockedList(in *friend.GetBlockedListRequest) (*friend.GetBlockedListResponse, error) {
	if in.UserId == 0 {
		return &friend.GetBlockedListResponse{Total: 0, Friends: []*friend.Friend{}}, nil
	}

	// 构建查询，只查询被我拉黑的好友
	query := l.svcCtx.DB.Model(&model.Friends{}).Where("user_id = ? AND status = ?", in.UserId, _const.FriendStatusBlocked)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("获取黑名单总数失败: %v", err)
		return nil, err // 直接返回错误，让框架处理
	}

	if total == 0 {
		return &friend.GetBlockedListResponse{Total: 0, Friends: []*friend.Friend{}}, nil
	}

	var blockedFriends []model.Friends
	// 应用分页和排序
	err := query.Order("update_at DESC").Find(&blockedFriends).Error
	if err != nil {
		l.Logger.Errorf("获取黑名单列表失败: %v", err)
		return nil, err
	}

	var result []*friend.Friend
	for _, f := range blockedFriends {
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

	return &friend.GetBlockedListResponse{
		Friends: result,
		Total:   total,
	}, nil
}
