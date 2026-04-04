package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ai-gateway/internal/domain/entity"
)

// TokenSyncJob token 同步job
type TokenSyncJob struct {
	db   *gorm.DB
	rdb  redis.UniversalClient
	stop chan struct{}
}

// NewTokenSyncJob 创建token同步job
func NewTokenSyncJob(db *gorm.DB, rdb redis.UniversalClient) *TokenSyncJob {
	return &TokenSyncJob{
		db:   db,
		rdb:  rdb,
		stop: make(chan struct{}, 1),
	}
}

// Start 启动job
func (j *TokenSyncJob) Start() {
	go func() {
		j.sync() // 进入执行一次

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("token-sync job is running...")
				j.sync()
				log.Println("token-sync job is finished")
			case <-j.stop:
				return
			}
		}
	}()
}

func (j *TokenSyncJob) sync() {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	ctx := context.Background()

	// 1. 同步月度Token到Redis（用于限流）
	var keys []entity.APIKey
	j.db.Where("status = 1 AND monthly_token_limit > 0").Find(&keys)
	for _, key := range keys {
		var totalTokens int64
		j.db.Model(&entity.RequestLog{}).
			Where("api_key_id = ? AND created_at >= ?", key.ID, monthStart).
			Select("COALESCE(SUM(tokens_used), 0)").
			Scan(&totalTokens)
		monthKey := fmt.Sprintf("monthly_tokens:%d", key.ID)
		j.rdb.Set(ctx, monthKey, totalTokens, 5*time.Minute)
	}

	// 2. 聚合今日Token使用量写入token_usages表
	type dailyStat struct {
		TenantID uint
		APIKeyID uint
		Tokens   int64
	}
	var stats []dailyStat
	j.db.Model(&entity.RequestLog{}).
		Select("tenant_id, api_key_id, COALESCE(SUM(tokens_used), 0) as tokens").
		Where("created_at >= ?", today).
		Group("tenant_id, api_key_id").
		Scan(&stats)

	for _, s := range stats {
		j.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "api_key_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"tokens"}),
		}).Create(&entity.TokenUsage{
			TenantID: s.TenantID,
			APIKeyID: s.APIKeyID,
			Tokens:   int(s.Tokens),
			Date:     today,
		})
	}
}

// Stop 停止job
func (j *TokenSyncJob) Stop() {
	close(j.stop)
}
