package model

import (
	"gorm.io/gorm"
	"time"
)

type Groups struct {
	Id             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Description    string         `gorm:"type:varchar(500);default:''" json:"description"`
	Avatar         string         `gorm:"type:varchar(255);default:''" json:"avatar"`
	OwnerId        int64          `gorm:"not null;index:idx_owner" json:"owner_id"`
	MemberCount    int64          `gorm:"default:1" json:"member_count"`
	MaxMemberCount int64          `gorm:"default:500" json:"max_member_count"`
	Status         int64          `gorm:"type:tinyint(4);default:1;index:idx_status" json:"status"` // 0 正常 1 禁用 2 解散
	CreateAt       int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt       int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type GroupMembers struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupId   int64          `gorm:"not null;index:idx_group_user,unique" json:"group_id"`
	UserId    int64          `gorm:"not null;index:idx_group_user,unique;index:idx_user_id" json:"user_id"`
	Role      int64          `gorm:"type:tinyint(4);default:3;index:idx_role" json:"role"` // 0 普通成员 1 群主 2 管理员
	Nickname  string         `gorm:"type:varchar(50);default:''" json:"nickname"`
	Status    int64          `gorm:"type:tinyint(4);default:1" json:"status"`      // 0 正常 1 禁言（标志位，非必须）
	MuteUntil time.Time      `gorm:"type:datetime;default:null" json:"mute_until"` // 新增字段，记录禁言结束时间
	JoinTime  int64          `gorm:"not null" json:"join_time"`
	UpdateAt  int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type JoinGroupApplications struct {
	Id         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserId int64          `gorm:"not null;index:idx_from_user" json:"from_user_id"`         // 申请人
	ToGroupId  int64          `gorm:"not null;index:idx_to_group" json:"to_group_id"`           // 目标群组
	Reason     string         `gorm:"type:varchar(200);default:''" json:"reason"`               // 申请理由
	InviterId  int64          `gorm:"not null;index:idx_inviter" json:"inviter_id"`             // 邀请人ID(默认为0，表示自己申请)
	Status     int8           `gorm:"type:tinyint(4);default:1;index:idx_status" json:"status"` // 0: 待处理 1: 同意 2: 拒绝
	CreateAt   int64          `gorm:"autoCreateTime" json:"create_time"`                        // 申请时间
	UpdateAt   int64          `gorm:"autoUpdateTime" json:"update_time"`                        // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
