package repo

import (
	"ai-gateway/internal/domain/entity"
)

// LogRepository 请求日志接口
type LogRepository interface {
	Create(entry *entity.RequestLog) error
	List(limit int, page int) ([]entity.RequestLog, int64, error)
}
