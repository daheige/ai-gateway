package repo

import (
	"ai-gateway/internal/domain/entity"
)

type StatsRepository interface {
	// GetOverview 获取总览统计
	GetOverview() (*entity.StatsOverview, error)

	// GetDailyStats 获取最近N天每日Token统计
	GetDailyStats(days int) ([]entity.DailyStat, error)

	GetKeyStats() ([]entity.KeyStat, error)
}
