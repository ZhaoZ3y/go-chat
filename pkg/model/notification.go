package model

import "gorm.io/gorm"

type Notifications struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64          `gorm:"not null;index:idx_user_id" json:"user_id"`
	Type      int8           `gorm:"type:tinyint(4);not null" json:"type"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Data      string         `gorm:"type:json" json:"data"`
	IsRead    bool           `gorm:"default:false;index:idx_is_read" json:"is_read"`
	CreateAt  int64          `gorm:"autoCreateTime;index:idx_create_time" json:"create_time"`
	UpdateAt  int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
