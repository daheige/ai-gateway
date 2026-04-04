package providers

import (
	"gorm.io/gorm"

	"ai-gateway/internal/domain/repo"
	"ai-gateway/internal/infras/persistence"
)

// Repositories 资源池
type Repositories struct {
	APIKeyRepo   repo.APIKeyRepository
	LogRepo      repo.LogRepository
	ProviderRepo repo.ProviderRepository
	StatsRepo    repo.StatsRepository
	TenantRepo   repo.TenantRepository
}

// NewRepositories 创建资源
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		APIKeyRepo:   persistence.NewAPIKeyRepo(db),
		LogRepo:      persistence.NewLogRepo(db),
		ProviderRepo: persistence.NewProviderRepo(db),
		StatsRepo:    persistence.NewStatsRepo(db),
		TenantRepo:   persistence.NewTenantRepo(db),
	}
}
