package entity

import (
	"time"
)

// TokenUsage token使用实体
type TokenUsage struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TenantID  uint      `gorm:"not null;uniqueIndex:idx_usage_unique" json:"tenant_id"`
	APIKeyID  uint      `gorm:"uniqueIndex:idx_usage_unique" json:"api_key_id"`
	Tokens    int       `json:"tokens"`
	Date      time.Time `gorm:"uniqueIndex:idx_usage_unique" json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (TokenUsage) TableName() string {
	return "token_usages"
}
