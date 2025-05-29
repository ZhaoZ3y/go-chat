package model

import "gorm.io/gorm"

type Files struct {
	Id           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename     string         `gorm:"type:varchar(255);not null" json:"filename"`      // MinIO 对象的基本名称 (例如：时间戳 + 扩展名)
	OriginalName string         `gorm:"type:varchar(255);not null" json:"original_name"` // 客户端上传的原始文件名
	FilePath     string         `gorm:"type:varchar(500);not null" json:"file_path"`     // 完整的 MinIO 对象键 (Object Key)
	FileUrl      string         `gorm:"type:varchar(500);not null" json:"file_url"`      // 访问 MinIO 文件的完整 URL
	FileType     string         `gorm:"type:varchar(50);not null;index:idx_type" json:"file_type"`
	FileSize     int64          `gorm:"not null" json:"file_size"`
	MimeType     string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	Hash         string         `gorm:"type:varchar(64);not null;index:idx_hash" json:"hash"`
	UserId       int64          `gorm:"not null;index:idx_user" json:"user_id"`
	Status       int8           `gorm:"type:tinyint(4);default:1" json:"status"` // 1:正常 2:已删除
	CreateAt     int64          `gorm:"autoCreateTime;index:idx_create_time" json:"create_time"`
	UpdateAt     int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
