package logic

import (
	_const "IM/pkg/const"
	"IM/pkg/model"
	"context"
	"gorm.io/gorm"
	"time"

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
	// 1. 参数校验
	if in.RequestId == 0 || in.UserId == 0 {
		return &friend.HandleFriendRequestResponse{Success: false, Message: "参数错误：申请ID和用户ID不能为空"}, nil
	}
	// 使用常量进行判断，更清晰
	if in.Action != _const.FriendRequestStatusAccepted && in.Action != _const.FriendRequestStatusRejected {
		return &friend.HandleFriendRequestResponse{Success: false, Message: "操作类型错误，仅支持同意或拒绝"}, nil
	}

	// 2. 查找待处理的好友申请
	var request model.FriendRequests
	// 只查找属于当前用户(to_user_id)且状态为“待处理”的申请
	err := l.svcCtx.DB.Where("id = ? AND to_user_id = ? AND status = ?", in.RequestId, in.UserId, _const.FriendRequestStatusPending).
		First(&request).Error
	if err != nil {
		// 区分是“没找到”还是“数据库错误”
		if err == gorm.ErrRecordNotFound {
			return &friend.HandleFriendRequestResponse{Success: false, Message: "好友申请不存在或已被处理"}, nil
		}
		l.Logger.Errorf("查询好友申请失败, RequestID: %d, Error: %v", in.RequestId, err)
		return &friend.HandleFriendRequestResponse{Success: false, Message: "系统错误，请稍后再试"}, nil
	}

	// 3. 开启事务处理
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. 更新申请状态，并标记为“已读”
	updateData := map[string]interface{}{
		"status":    in.Action,
		"is_read":   true, // 无论同意或拒绝，都标记为已读
		"update_at": time.Now().Unix(),
	}
	if err := tx.Model(&request).Updates(updateData).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新好友申请状态失败: %v", err)
		return &friend.HandleFriendRequestResponse{Success: false, Message: "处理申请失败"}, nil
	}

	// 5. 如果同意申请，创建或更新双向好友关系
	if in.Action == _const.FriendRequestStatusAccepted {
		// 使用 FirstOrCreate 保证幂等性，如果之前是好友后被删除，则可以恢复关系
		// a. 创建 A -> B 的关系
		friendAB := model.Friends{UserId: request.FromUserId, FriendId: request.ToUserId, Status: _const.FriendStatusNormal}
		if err := tx.Where(model.Friends{UserId: request.FromUserId, FriendId: request.ToUserId}).Assign(friendAB).FirstOrCreate(&model.Friends{}).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建好友关系 A->B 失败: %v", err)
			return &friend.HandleFriendRequestResponse{Success: false, Message: "建立好友关系失败"}, nil
		}

		// b. 创建 B -> A 的关系
		friendBA := model.Friends{UserId: request.ToUserId, FriendId: request.FromUserId, Status: _const.FriendStatusNormal}
		if err := tx.Where(model.Friends{UserId: request.ToUserId, FriendId: request.FromUserId}).Assign(friendBA).FirstOrCreate(&model.Friends{}).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建好友关系 B->A 失败: %v", err)
			return &friend.HandleFriendRequestResponse{Success: false, Message: "建立好友关系失败"}, nil
		}
	}

	// 6. 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("处理好友申请事务提交失败: %v", err)
		return &friend.HandleFriendRequestResponse{Success: false, Message: "处理申请失败"}, nil
	}

	// 7. 构造成功的响应
	message := "已拒绝好友申请"
	if in.Action == _const.FriendRequestStatusAccepted {
		message = "已同意好友申请，你们现在是好友了"
	}

	// 【重要修复】填充响应中的 RequestInfo 字段，供API层使用
	return &friend.HandleFriendRequestResponse{
		Success: true,
		Message: message,
		RequestInfo: &friend.FriendRequest{
			Id:         request.Id,
			FromUserId: request.FromUserId,
			ToUserId:   request.ToUserId,
			Message:    request.Message,
			Status:     in.Action, // 返回最新的状态
			CreateAt:   request.CreateAt,
			UpdateAt:   time.Now().Unix(),
		},
	}, nil
}
