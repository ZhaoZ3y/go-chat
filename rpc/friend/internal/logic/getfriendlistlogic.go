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

// 获取好友列表
func (l *GetFriendListLogic) GetFriendList(in *friend.GetFriendListRequest) (*friend.GetFriendListResponse, error) {
	// 1. 参数验证
	if in.UserId == 0 {
		return &friend.GetFriendListResponse{}, nil
	}

	// 2. 从数据库获取好友基础列表
	var friends []model.Friends
	query := l.svcCtx.DB.Model(&model.Friends{}).Where("user_id = ? AND status = 1", in.UserId)
	if err := query.Find(&friends).Error; err != nil {
		l.Logger.Errorf("获取好友列表失败: %v", err)
		return nil, err // 建议返回错误，而不是空的响应
	}

	if len(friends) == 0 {
		return &friend.GetFriendListResponse{Total: 0, Friends: []*friend.Friend{}}, nil
	}

	// 3. 准备批量获取在线状态
	friendIDs := make([]int64, len(friends))
	for i, f := range friends {
		friendIDs[i] = f.FriendId
	}

	statusMap, err := l.svcCtx.UserStatusSvc.GetBatchUserStatus(friendIDs)
	if err != nil {
		l.Logger.Errorf("批量获取好友在线状态失败: %v", err)
		statusMap = make(map[int64]*chat_service.UserStatus) // 创建一个空map防止下面代码panic
	}

	aggregatedList := make([]friendWithStatus, 0, len(friends))
	pinyinArgs := pinyin.NewArgs() // 创建拼音转换参数

	for _, f := range friends {
		status, ok := statusMap[f.FriendId]
		onlineStatus := int64(0)
		if ok && status.Status == 1 {
			onlineStatus = 1
		}

		// 生成排序键 (首字母)
		remark := strings.TrimSpace(f.Remark)
		if remark == "" {
			remark = "#" // 如果备注为空，则归入'#'组
		}
		pinyinSlice := pinyin.Pinyin(remark, pinyinArgs)
		firstChar := pinyinSlice[0][0][0]
		sortKey := "#"
		if unicode.IsLetter(rune(firstChar)) {
			sortKey = strings.ToUpper(string(firstChar))
		}

		aggregatedList = append(aggregatedList, friendWithStatus{
			FriendData:   f,
			OnlineStatus: onlineStatus,
			SortKey:      sortKey,
		})
	}

	// 6. 执行自定义排序
	sort.Slice(aggregatedList, func(i, j int) bool {
		// 主排序规则：按首字母 A-Z 排序
		if aggregatedList[i].SortKey != aggregatedList[j].SortKey {
			// 将'#'排在最后
			if aggregatedList[i].SortKey == "#" {
				return false
			}
			if aggregatedList[j].SortKey == "#" {
				return true
			}
			return aggregatedList[i].SortKey < aggregatedList[j].SortKey
		}

		// 次排序规则：首字母相同时，在线的排在前面
		return aggregatedList[i].OnlineStatus > aggregatedList[j].OnlineStatus
	})

	// 7. 转换为最终响应格式
	result := make([]*friend.Friend, len(aggregatedList))
	for i, item := range aggregatedList {
		f := item.FriendData
		result[i] = &friend.Friend{
			Id:           f.Id,
			UserId:       f.UserId,
			FriendId:     f.FriendId,
			Remark:       f.Remark,
			Status:       int32(f.Status),
			CreateAt:     f.CreateAt,
			UpdateAt:     f.UpdateAt,
			OnlineStatus: item.OnlineStatus, // 填充在线状态
		}
	}

	return &friend.GetFriendListResponse{
		Friends: result,
		Total:   int64(len(friends)),
	}, nil
}
