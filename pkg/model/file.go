package model

import "gorm.io/gorm"

type Files struct {
	Id           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename     string         `gorm:"type:varchar(255);not null" json:"filename"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FilePath     string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileUrl      string         `gorm:"type:varchar(500);not null" json:"file_url"`
	FileType     string         `gorm:"type:varchar(50);not null;index:idx_type" json:"file_type"`
	FileSize     int64          `gorm:"not null" json:"file_size"`
	MimeType     string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	Hash         string         `gorm:"type:varchar(64);not null;index:idx_hash" json:"hash"`
	UserId       int64          `gorm:"not null;index:idx_user" json:"user_id"`
	Status       int8           `gorm:"type:tinyint(4);default:1" json:"status"`
	CreateAt     int64          `gorm:"autoCreateTime;index:idx_create_time" json:"create_time"`
	UpdateAt     int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type FileUploads struct {
	Id             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UploadId       string         `gorm:"type:varchar(64);not null;uniqueIndex:idx_upload_id" json:"upload_id"`
	Filename       string         `gorm:"type:varchar(255);not null" json:"filename"`
	FileSize       int64          `gorm:"not null" json:"file_size"`
	ChunkSize      int            `gorm:"not null" json:"chunk_size"`
	TotalChunks    int            `gorm:"not null" json:"total_chunks"`
	UploadedChunks string         `gorm:"type:text" json:"uploaded_chunks"`
	UserId         int64          `gorm:"not null;index:idx_user" json:"user_id"`
	Status         int8           `gorm:"type:tinyint(4);default:1;index:idx_status" json:"status"`
	CreateAt       int64          `gorm:"autoCreateTime" json:"create_time"`
	UpdateAt       int64          `gorm:"autoUpdateTime" json:"update_time"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
