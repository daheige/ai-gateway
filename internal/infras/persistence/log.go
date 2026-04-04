package persistence

import (
	"log"

	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

type LogRepoImpl struct {
	db *gorm.DB
}

// NewLogRepo 创建日志接口实例
func NewLogRepo(db *gorm.DB) repo.LogRepository {
	return &LogRepoImpl{db: db}
}

// Create 记录操作日志
func (l *LogRepoImpl) Create(entry *entity.RequestLog) error {
	return l.db.Create(entry).Error
}

// List 日志列表
func (l *LogRepoImpl) List(limit int, page int) ([]entity.RequestLog, int64, error) {
	var total int64
	err := l.db.Model(&entity.RequestLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []entity.RequestLog{}, 0, nil
	}

	log.Println("total:", total)
	var logs []entity.RequestLog
	err = l.db.Order("id desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&logs).Error
	log.Println("logs:", logs, "err:", err)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
