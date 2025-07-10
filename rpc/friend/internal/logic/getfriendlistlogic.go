package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/chat_service"
	"context"
	"github.com/mozillazg/go-pinyin"
	"sort"
	"strings"
	"unicode"

	"IM/rpc/friend/friend"
	"IM/rpc/friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type friendWithStatus struct {
	FriendData   model.Friends
	FriendUser   model.User
	OnlineStatus int64  // 0: 离线, 1: 在线
	SortKey      string // 用于排序的首字母
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *friend.GetFriendListRequest) (*friend.GetFriendListResponse, error) {
	if in.UserId == 0 {
		return &friend.GetFriendListResponse{}, nil
	}

	// 1. 先查好友关系列表
	var friends []model.Friends
	err := l.svcCtx.DB.Where("user_id = ? AND status = 1", in.UserId).Find(&friends).Error
	if err != nil {
		l.Logger.Errorf("获取好友列表失败: %v", err)
		return nil, err
	}
	if len(friends) == 0 {
		return &friend.GetFriendListResponse{Total: 0, Friends: []*friend.Friend{}}, nil
	}

	friendIDs := make([]int64, len(friends))
	for i, f := range friends {
		friendIDs[i] = f.FriendId
	}

	// 2. 批量获取好友用户信息（头像，昵称）
	var users []model.User
	err = l.svcCtx.DB.Where("id IN ?", friendIDs).Find(&users).Error
	if err != nil {
		l.Logger.Errorf("批量获取好友用户信息失败: %v", err)
		return nil, err
	}
	userMap := make(map[int64]model.User, len(users))
	for _, u := range users {
		userMap[u.Id] = u
	}

	// 3. 批量获取好友在线状态
	statusMap, err := l.svcCtx.UserStatusSvc.GetBatchUserStatus(friendIDs)
	if err != nil {
		l.Logger.Errorf("批量获取好友在线状态失败: %v", err)
		statusMap = make(map[int64]*chat_service.UserStatus)
	}

	pinyinArgs := pinyin.NewArgs()
	aggregatedList := make([]friendWithStatus, 0, len(friends))

	for _, f := range friends {
		u, ok := userMap[f.FriendId]
		if !ok {
			continue
		}

		status, ok := statusMap[f.FriendId]
		onlineStatus := int64(0)
		if ok && status.Status == 1 {
			onlineStatus = 1
		}

		// 备注为空时用昵称代替
		displayName := strings.TrimSpace(f.Remark)
		if displayName == "" {
			displayName = strings.TrimSpace(u.Nickname)
		}
		if displayName == "" {
			displayName = "#" // 兜底
		}

		// 生成排序键：拼音首字母大写
		pinyinSlice := pinyin.Pinyin(displayName, pinyinArgs)
		sortKey := "#"
		if len(pinyinSlice) > 0 && len(pinyinSlice[0]) > 0 && len(pinyinSlice[0][0]) > 0 {
			firstChar := pinyinSlice[0][0][0]
			if unicode.IsLetter(rune(firstChar)) {
				sortKey = strings.ToUpper(string(firstChar))
			}
		}

		aggregatedList = append(aggregatedList, friendWithStatus{
			FriendData:   f,
			FriendUser:   u,
			OnlineStatus: onlineStatus,
			SortKey:      sortKey,
		})
	}

	// 排序：首字母升序，'#'排后面；在线排前面
	sort.Slice(aggregatedList, func(i, j int) bool {
		if aggregatedList[i].SortKey != aggregatedList[j].SortKey {
			if aggregatedList[i].SortKey == "#" {
				return false
			}
			if aggregatedList[j].SortKey == "#" {
				return true
			}
			return aggregatedList[i].SortKey < aggregatedList[j].SortKey
		}
		return aggregatedList[i].OnlineStatus > aggregatedList[j].OnlineStatus
	})

	// 组装返回
	result := make([]*friend.Friend, len(aggregatedList))
	for i, item := range aggregatedList {
		f := item.FriendData
		u := item.FriendUser
		result[i] = &friend.Friend{
			Id:           f.Id,
			UserId:       f.UserId,
			FriendId:     f.FriendId,
			Remark:       f.Remark,
			Status:       int32(f.Status),
			CreateAt:     f.CreateAt,
			UpdateAt:     f.UpdateAt,
			OnlineStatus: item.OnlineStatus,
			Nickname:     u.Nickname,
			Avatar:       u.Avatar,
		}
	}

	return &friend.GetFriendListResponse{
		Friends: result,
		Total:   int64(len(result)),
	}, nil
}
