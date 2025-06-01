package database

import (
	"IM/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(c config.Config) (db *gorm.DB, err error) {
	// 这里需要根据实际情况配置数据库连接字符串
	db, err = gorm.Open(mysql.Open(c.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
