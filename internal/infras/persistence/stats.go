package persistence

import (
	"time"

	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

// StatsRepoImpl stats repo 接口实现
type StatsRepoImpl struct {
	db *gorm.DB
}

var _ repo.StatsRepository = (*StatsRepoImpl)(nil)

// NewStatsRepo 创建stats repo接口实例
func NewStatsRepo(db *gorm.DB) repo.StatsRepository {
	return &StatsRepoImpl{db: db}
}

// GetOverview 获取总览统计
func (s *StatsRepoImpl) GetOverview() (*entity.StatsOverview, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var overview entity.StatsOverview
	s.db.Model(&entity.RequestLog{}).
		Where("created_at >= ?", today).
		Select("COALESCE(SUM(tokens_used), 0)").
		Scan(&overview.TotalTokensToday)

	s.db.Model(&entity.RequestLog{}).
		Where("created_at >= ?", monthStart).
		Select("COALESCE(SUM(tokens_used), 0)").
		Scan(&overview.TotalTokensMonth)

	s.db.Model(&entity.RequestLog{}).Count(&overview.TotalRequests)

	s.db.Model(&entity.RequestLog{}).
		Where("created_at >= ?", today).
		Count(&overview.TotalRequestsDay)

	s.db.Model(&entity.APIKey{}).Where("status = 1").Count(&overview.ActiveKeys)
	s.db.Model(&entity.Provider{}).Where("status = 1").Count(&overview.ActiveProviders)

	return &overview, nil
}

// GetDailyStats 获取最近N天每日Token统计
func (s *StatsRepoImpl) GetDailyStats(days int) ([]entity.DailyStat, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	var stats []entity.DailyStat
	err := s.db.Model(&entity.TokenUsage{}).
		Select("DATE_FORMAT(date, '%Y-%m-%d') as date, SUM(tokens) as tokens").
		Where("date >= ?", since).
		Group("DATE_FORMAT(date, '%Y-%m-%d')").
		Order("date").
		Scan(&stats).Error

	return stats, err
}

// GetKeyStats 获取当月各Key的Token使用排行
func (s *StatsRepoImpl) GetKeyStats() ([]entity.KeyStat, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var stats []entity.KeyStat
	err := s.db.Model(&entity.RequestLog{}).
		Select("request_logs.api_key_id, api_keys.key_prefix, COALESCE(SUM(request_logs.tokens_used), 0) as tokens").
		Joins("LEFT JOIN api_keys ON api_keys.id = request_logs.api_key_id").
		Where("request_logs.created_at >= ?", monthStart).
		Group("request_logs.api_key_id, api_keys.key_prefix").
		Order("tokens DESC").
		Limit(20).
		Scan(&stats).Error

	return stats, err
}
