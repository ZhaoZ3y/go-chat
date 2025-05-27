package model

import "gorm.io/gorm"

type User struct {
	Id        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex:idx_username;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex:idx_email;not null" json:"email"`
	Nickname  string         `gorm:"type:varchar(50);default:''" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255);default:''" json:"avatar"`
	Phone     string         `gorm:"type:varchar(20);default:'';index:idx_phone" json:"phone"`
	Status    int8           `gorm:"type:tinyint(4);default:1" json:"status"`
	CreateAt  int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt  int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
