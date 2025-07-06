package model

import (
	"gorm.io/gorm"
)

type FileRecord struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID      string         `gorm:"uniqueIndex;size:64;not null" json:"file_id"`
	FileName    string         `gorm:"size:255;not null" json:"file_name"`
	FileType    string         `gorm:"size:32;not null;index" json:"file_type"`
	FileSize    int64          `gorm:"not null" json:"file_size"`
	ContentType string         `gorm:"size:100;not null" json:"content_type"`
	ETag        string         `gorm:"size:128" json:"etag"`
	ObjectName  string         `gorm:"size:255;not null" json:"object_name"`
	UserID      int64          `gorm:"not null;index" json:"user_id"`
	Status      int            `gorm:"default:1" json:"status"` // 1:正常 2:已删除
	CreateAt    int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt    int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段
	ExpireAt    int64          `gorm:"index" json:"expire_at"`
}
