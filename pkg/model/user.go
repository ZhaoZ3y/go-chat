package model

type User struct {
	Id          int64  `gorm:"primaryKey;autoIncrement;comment:用户ID" json:"id"`
	Username    string `gorm:"unique;not null;comment:用户名" json:"username"`
	Password    string `gorm:"not null;comment:密码" json:"-"`
	Nickname    string `gorm:"not null;comment:昵称" json:"nickname"`
	Avatar      string `gorm:"comment:头像" json:"avatar"`
	Email       string `gorm:"unique;comment:邮箱" json:"email"`
	Phone       string `gorm:"unique;comment:手机号" json:"phone"`
	Status      int32  `gorm:"not null;default:1;comment:状态 1-正常 2-禁用" json:"status"`
	CreateAt    int64  `gorm:"autoCreateTime;comment:创建时间" json:"create_at"`
	UpdateAt    int64  `gorm:"autoUpdateTime;comment:更新时间" json:"update_at"`
	LastLoginAt int64  `gorm:"comment:最后登录时间" json:"last_login_at"`
	DeletedAt   int64  `gorm:"comment:删除时间" json:"deleted_at"`
}
