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
	RateLimitPerMin   int            `gorm:"default:60" json:"rate_limit_per_min"`
	MonthlyTokenLimit int            `gorm:"default:0" json:"monthly_token_limit"`
	TotalTokenQuota   int            `gorm:"default:0" json:"total_token_quota"`
	TokensConsumed    int            `gorm:"default:0" json:"tokens_consumed"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (APIKey) TableName() string {
	return "api_keys"
}
