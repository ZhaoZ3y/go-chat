package model

import (
	"gorm.io/gorm"
)

type Messages struct {
	Id         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserId int64          `gorm:"not null;index:idx_from_user" json:"from_user_id"`
	ToUserId   int64          `gorm:"default:0;index:idx_to_user" json:"to_user_id"`
	GroupId    int64          `gorm:"default:0;index:idx_group" json:"group_id"`
	Type       int64          `gorm:"type:tinyint(4);default:0" json:"type"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Extra      string         `gorm:"type:text" json:"extra"`
	ChatType   int64          `gorm:"type:tinyint(4);default:0;index:idx_chat_type" json:"chat_type"`
	Status     int8           `gorm:"type:tinyint(4);default:0;index:idx_status" json:"status"` // 0:正常 1:已撤回 2.删除
	CreateAt   int64          `gorm:"autoCreateTime;index:idx_create_time" json:"create_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// MessageReadReceipts 已读回执表
type MessageReadReceipts struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageId int64          `gorm:"not null;uniqueIndex:uq_msg_user" json:"message_id"`
	UserId    int64          `gorm:"not null;uniqueIndex:uq_msg_user" json:"user_id"`
	GroupId   int64          `gorm:"default:0;index" json:"group_id"`
	ChatType  int8           `gorm:"type:tinyint(4);default:0" json:"chat_type"`
	ReadAt    int64          `gorm:"autoCreateTime" json:"read_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MessageUserStates struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageId int64          `gorm:"not null;uniqueIndex:uq_msg_user" json:"message_id"`
	UserId    int64          `gorm:"not null;uniqueIndex:uq_msg_user" json:"user_id"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Conversations struct {
	Id              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId          int64          `gorm:"not null;uniqueIndex:uq_user_target_type" json:"user_id"`
	TargetId        int64          `gorm:"not null;uniqueIndex:uq_user_target_type" json:"target_id"`
	Type            int8           `gorm:"type:tinyint;default:0;uniqueIndex:uq_user_target_type" json:"type"`
	LastMessageID   int64          `gorm:"default:0" json:"last_message_id"` // 存储最后一条消息的ID
	LastMessageTime int64          `gorm:"default:0" json:"last_message_time"`
	LastMessage     string         `gorm:"type:varchar(255);default:''" json:"last_message"` // 最后一条消息内容
	UnreadCount     int            `gorm:"default:0" json:"unread_count"`
	IsPinned        bool           `gorm:"default:false;index" json:"is_pinned"` // 是否置顶
	CreateAt        int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt        int64          `gorm:"autoUpdateTime;index" json:"update_time"` // 更新时间索引很常用
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type OfflineMessages struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64          `gorm:"not null;index:idx_user" json:"user_id"`
	MessageId int64          `gorm:"not null;index:idx_message" json:"message_id"`
	Type      int8           `gorm:"type:tinyint(4);default:0" json:"type"` // 0:普通消息 1:系统通知
	Content   string         `gorm:"type:text" json:"content"`
	Status    int8           `gorm:"type:tinyint(4);default:0" json:"status"` // 0:未推送 1:已推送
	CreateAt  int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt  int64          `gorm:"autoUpdateTime;index:idx_update_time" json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
