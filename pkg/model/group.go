package model

import "gorm.io/gorm"

type Groups struct {
	Id             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Description    string         `gorm:"type:varchar(500);default:''" json:"description"`
	Avatar         string         `gorm:"type:varchar(255);default:''" json:"avatar"`
	OwnerId        int64          `gorm:"not null;index:idx_owner" json:"owner_id"`
	MemberCount    int            `gorm:"default:1" json:"member_count"`
	MaxMemberCount int            `gorm:"default:500" json:"max_member_count"`
	Status         int8           `gorm:"type:tinyint(4);default:1;index:idx_status" json:"status"`
	CreateAt       int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt       int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type GroupMembers struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupId   int64          `gorm:"not null;index:idx_group_user,unique" json:"group_id"`
	UserId    int64          `gorm:"not null;index:idx_group_user,unique;index:idx_user_id" json:"user_id"`
	Role      int8           `gorm:"type:tinyint(4);default:3;index:idx_role" json:"role"`
	Nickname  string         `gorm:"type:varchar(50);default:''" json:"nickname"`
	Status    int8           `gorm:"type:tinyint(4);default:1" json:"status"`
	JoinTime  int64          `gorm:"not null" json:"join_time"`
	UpdateAt  int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
