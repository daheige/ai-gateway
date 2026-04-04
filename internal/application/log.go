package application

import (
	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

// LogService 日志服务
type LogService struct {
	logRepo repo.LogRepository
}

// NewLogService 创建日志服务实例
func NewLogService(logRepo repo.LogRepository) *LogService {
	return &LogService{logRepo: logRepo}
}

// Create 创建日志
func (l *LogService) Create(entry *entity.RequestLog) error {
	// todo 这里日志插入，可以将其丢入redis list或kafka mq中，然后通过异步job来消费
	return l.logRepo.Create(entry)
}

// 日志列表
func (l *LogService) List(limit int, page int) ([]entity.RequestLog, int64, error) {
	return l.logRepo.List(limit, page)
}
