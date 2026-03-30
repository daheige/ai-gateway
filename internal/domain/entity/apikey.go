package entity

import (
	"time"

	"gorm.io/gorm"
)

// APIKey 虚拟apikey配置
type APIKey struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	TenantID          uint           `gorm:"not null;index" json:"tenant_id"`
	KeyHash           string         `gorm:"size:64;uniqueIndex;not null" json:"-"`
	KeyPrefix         string         `gorm:"size:20" json:"key_prefix"`
	Name              string         `gorm:"size:100" json:"name"`
	ProviderID        uint           `gorm:"index" json:"provider_id"`
	Status            int            `gorm:"default:1" json:"status"`
	RateLimitPerSec   int            `gorm:"default:10" json:"rate_limit_per_sec"`
	RateLimit         int            `gorm:"default:60" json:"rate_limit"`
	MonthlyTokenLimit int            `gorm:"default:0" json:"monthly_token_limit"`
	TotalTokenQuota   int            `gorm:"default:0" json:"total_token_quota"`
	TokensConsumed    int            `gorm:"default:0" json:"tokens_consumed"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// RequestLog 请求日志实体
type RequestLog struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	TenantID     uint      `gorm:"not null;index" json:"tenant_id"`
	APIKeyID     uint      `gorm:"index" json:"api_key_id"`
	ProviderID   uint      `gorm:"index" json:"provider_id"`
	Model        string    `gorm:"size:50" json:"model"`
	TokensUsed   int       `gorm:"default:0" json:"tokens_used"`
	PromptTokens int       `gorm:"default:0" json:"prompt_tokens"`
	CompTokens   int       `gorm:"default:0" json:"comp_tokens"`
	Status       int       `json:"status"`
	Latency      int       `json:"latency"`
	IP           string    `gorm:"size:50" json:"ip"`
	CreatedAt    time.Time `json:"created_at"`
}

// TokenUsage token使用实体
type TokenUsage struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TenantID  uint      `gorm:"not null;uniqueIndex:idx_usage_unique" json:"tenant_id"`
	APIKeyID  uint      `gorm:"uniqueIndex:idx_usage_unique" json:"api_key_id"`
	Tokens    int       `json:"tokens"`
	Date      time.Time `gorm:"uniqueIndex:idx_usage_unique" json:"date"`
	CreatedAt time.Time `json:"created_at"`
}
