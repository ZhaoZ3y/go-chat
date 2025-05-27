package model

import "gorm.io/gorm"

type Messages struct {
	Id         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserId int64          `gorm:"not null;index:idx_from_user" json:"from_user_id"`
	ToUserId   int64          `gorm:"default:0;index:idx_to_user" json:"to_user_id"`
	GroupId    int64          `gorm:"default:0;index:idx_group" json:"group_id"`
	Type       int8           `gorm:"type:tinyint(4);default:0" json:"type"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Extra      string         `gorm:"type:text" json:"extra"`
	ChatType   int8           `gorm:"type:tinyint(4);default:0;index:idx_chat_type" json:"chat_type"`
	Status     int8           `gorm:"type:tinyint(4);default:1" json:"status"`
	CreateAt   int64          `gorm:"autoCreateTime;index:idx_create_time" json:"create_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type Conversations struct {
	Id              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId          int64          `gorm:"not null;index:idx_user_target_type,unique" json:"user_id"`
	TargetId        int64          `gorm:"not null;index:idx_user_target_type,unique" json:"target_id"`
	Type            int8           `gorm:"type:tinyint(4);default:0;index:idx_user_target_type,unique" json:"type"`
	LastMessage     string         `gorm:"type:text" json:"last_message"`
	LastMessageTime int64          `gorm:"default:0" json:"last_message_time"`
	UnreadCount     int            `gorm:"default:0" json:"unread_count"`
	CreateAt        int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt        int64          `gorm:"autoUpdateTime;index:idx_update_time" json:"update_time"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
