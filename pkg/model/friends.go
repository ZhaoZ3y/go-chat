package model

import "gorm.io/gorm"

type Friends struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64          `gorm:"not null;index:idx_user_friend,unique" json:"user_id"`
	FriendId  int64          `gorm:"not null;index:idx_user_friend,unique;index:idx_friend_id" json:"friend_id"`
	Remark    string         `gorm:"type:varchar(50);default:''" json:"remark"`
	Status    int8           `gorm:"type:tinyint(4);default:1" json:"status"`
	CreateAt  int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt  int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type FriendRequests struct {
	Id         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserId int64          `gorm:"not null;index:idx_from_user" json:"from_user_id"`
	ToUserId   int64          `gorm:"not null;index:idx_to_user" json:"to_user_id"`
	Message    string         `gorm:"type:varchar(200);default:''" json:"message"`
	Status     int8           `gorm:"type:tinyint(4);default:1;index:idx_status" json:"status"`
	CreateAt   int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt   int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
