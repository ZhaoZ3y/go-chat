package model

type GroupNotification struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Type         int64  `gorm:"type:tinyint(4);not null;index:idx_type" json:"type"`  // 通知类型
	GroupId      int64  `gorm:"not null;index:idx_group" json:"group_id"`             // 群组ID
	OperatorId   int64  `gorm:"not null;index:idx_operator" json:"operator_id"`       // 操作者ID
	TargetUserId int64  `gorm:"not null;index:idx_target_user" json:"target_user_id"` // 目标用户ID
	Message      string `gorm:"type:varchar(500);default:''" json:"message"`          // 通知内容
	IsRead       bool   `gorm:"default:false;index:idx_is_read" json:"is_read"`       // 是否已读
	CreateAt     int64  `gorm:"autoCreateTime" json:"create_time"`                    // 创建时间
	UpdateAt     int64  `gorm:"autoUpdateTime" json:"update_time"`                    // 更新时间
	DeletedAt    int64  `gorm:"index" json:"deleted_at"`                              // 删除时间
}
