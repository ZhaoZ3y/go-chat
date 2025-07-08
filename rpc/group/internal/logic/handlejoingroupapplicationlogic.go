package logic

import (
	"IM/pkg/model" // 假设您的数据库模型在这个包下
	"IM/pkg/mq"
	_const "IM/pkg/utils/const"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleJoinGroupApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleJoinGroupApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleJoinGroupApplicationLogic {
	return &HandleJoinGroupApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// HandleJoinGroupApplication 处理加入群组申请
func (l *HandleJoinGroupApplicationLogic) HandleJoinGroupApplication(in *group.HandleJoinGroupApplicationRequest) (*group.HandleJoinGroupApplicationResponse, error) {
	var application model.JoinGroupApplications
	err := l.svcCtx.DB.Where("id = ? AND status = ?", in.ApplicationId, group.ApplicationStatus_PENDING).
		First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("HandleJoinGroupApplication: 申请不存在或已被处理, id: %d", in.ApplicationId)
			return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "申请不存在或已被处理"}, nil
		}
		l.Logger.Errorf("查询入群申请失败, ApplicationID: %d, Error: %v", in.ApplicationId, err)
		return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "系统错误，查询申请失败"}, nil
	}

	// 2. 验证操作者权限 (必须是群主或管理员)
	var operatorMember model.GroupMembers // 假设这是您的群成员模型
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ?", application.ToGroupId, in.OperatorId).First(&operatorMember).Error
	if err != nil || (operatorMember.Role != int64(group.MemberRole_ROLE_OWNER) && operatorMember.Role != int64(group.MemberRole_ROLE_ADMIN)) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("HandleJoinGroupApplication: 操作者 %d 不是群 %d 的成员", in.OperatorId, application.ToGroupId)
		} else {
			logx.Errorf("HandleJoinGroupApplication: 操作者 %d 在群 %d 中没有权限, role: %d", in.OperatorId, application.ToGroupId, operatorMember.Role)
		}
		return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "您没有管理员权限，无法处理该申请"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	newStatus := group.ApplicationStatus_REJECTED
	if in.Approve {
		newStatus = group.ApplicationStatus_APPROVED
	}
	updateData := map[string]interface{}{
		"status":    int8(newStatus),
		"update_at": time.Now().Unix(),
	}
	if err := tx.Model(&application).Updates(updateData).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新入群申请状态失败: %v", err)
		return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "处理申请失败"}, nil
	}

	// 定义一个变量，用于在事务提交后引用新创建的欢迎消息
	var welcomeMessage *model.Messages

	// 5. 如果是同意申请，则执行核心逻辑
	if in.Approve {
		var existingMemberCount int64
		tx.Model(&model.GroupMembers{}).Where("group_id = ? AND user_id = ?", application.ToGroupId, application.FromUserId).Count(&existingMemberCount)

		if existingMemberCount == 0 {
			var applicantUser model.User
			if err := tx.First(&applicantUser, application.FromUserId).Error; err != nil {
				tx.Rollback()
				logx.Errorf("HandleJoinGroupApplication: 获取申请人信息失败, user_id: %d, error: %v", application.FromUserId, err)
				return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "获取申请人信息失败"}, nil
			}

			newMember := &model.GroupMembers{
				GroupId:  application.ToGroupId,
				UserId:   application.FromUserId,
				Nickname: applicantUser.Nickname,
				Role:     int64(group.MemberRole_ROLE_MEMBER),
				Status:   int64(group.MemberStatus_MEMBER_STATUS_NORMAL),
				JoinTime: time.Now().Unix(),
			}
			if err := tx.Create(newMember).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("创建群组成员记录失败: %v", err)
				return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "添加群组成员失败"}, nil
			}

			welcomeMessage = &model.Messages{
				FromUserId: _const.System,
				GroupId:    application.ToGroupId,
				Content:    fmt.Sprintf("欢迎 %s 加入群聊", newMember.Nickname),
				ChatType:   _const.ChatTypeGroup,
				Type:       _const.MsgTypeSystem,
			}
			if err := tx.Create(welcomeMessage).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("创建欢迎消息失败: %v", err)
				return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "创建欢迎消息失败"}, nil
			}

			// 5.3 更新群组的成员总数
			if err := tx.Model(&model.Groups{}).Where("id = ?", application.ToGroupId).UpdateColumn("member_count", gorm.Expr("member_count + 1")).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("更新群成员数量失败: %v", err)
				return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "更新群成员数量失败"}, nil
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("处理入群申请事务提交失败: %v", err)
		return &group.HandleJoinGroupApplicationResponse{Success: false, Message: "处理申请失败"}, nil
	}

	if welcomeMessage != nil {
		go l.publishWelcomeMessage(welcomeMessage.Id)
	}

	message := "已拒绝该用户的入群申请"
	if in.Approve {
		message = "已同意该用户的入群申请"
	}

	return &group.HandleJoinGroupApplicationResponse{
		Success: true,
		Message: message,
	}, nil
}

// publishWelcomeMessage 是一个辅助函数，用于异步地将新消息事件发布到 Kafka
func (l *HandleJoinGroupApplicationLogic) publishWelcomeMessage(messageId int64) {
	var finalMessage model.Messages
	if err := l.svcCtx.DB.First(&finalMessage, messageId).Error; err != nil {
		l.Logger.Errorf("发布欢迎消息前查询消息体失败: %v", err)
		return
	}

	event := &mq.MessageEvent{
		Type:        mq.EventNewMessage,
		MessageID:   finalMessage.Id,
		FromUserID:  finalMessage.FromUserId,
		GroupID:     finalMessage.GroupId,
		ChatType:    finalMessage.ChatType,
		MessageType: finalMessage.Type,
		Content:     finalMessage.Content,
		CreateAt:    finalMessage.CreateAt,
	}

	// 将事件发送到消息队列
	err := l.svcCtx.Kafka.SendMessage(mq.TopicMessage, event)
	if err != nil {
		l.Logger.Errorf("发布欢迎消息到 Kafka 失败: %v", err)
	}
}
