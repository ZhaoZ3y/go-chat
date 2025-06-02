package database

import (
	"IM/pkg/config"
	"IM/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
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

	if err := db.AutoMigrate(&model.Notifications{}); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db, nil
}
