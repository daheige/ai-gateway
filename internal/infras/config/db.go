package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化db
func InitDB(conf DatabaseConfig) *gorm.DB {
	// 数据库连接
	dbConfig := &gorm.Config{}
	if conf.ShowSQL {
		dbConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(conf.DSN), dbConfig)
	if err != nil {
		log.Fatalln("failed to open db err:", err)
	}

	return db
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("get db err:", err)
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatal("close db err:", err)
	}
}
