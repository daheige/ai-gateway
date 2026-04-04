package application

import (
	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

type StatsService struct {
	statsRepo repo.StatsRepository
}

// NewStatsService 创建stats服务实例
func NewStatsService(statsRepo repo.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

// GetOverview 获取总览统计
func (s *StatsService) GetOverview() (*entity.StatsOverview, error) {
	statsOverview, err := s.statsRepo.GetOverview()
	return statsOverview, err
}

// GetDailyStats 获取最近N天每日Token统计
func (s *StatsService) GetDailyStats(days int) ([]entity.DailyStat, error) {
	stats, err := s.statsRepo.GetDailyStats(days)
	return stats, err
}

// GetKeyStats 获取当月各Key的Token使用排行
func (s *StatsService) GetKeyStats() ([]entity.KeyStat, error) {
	stats, err := s.statsRepo.GetKeyStats()
	return stats, err
}
