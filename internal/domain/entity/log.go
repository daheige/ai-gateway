package entity

import (
	"time"
)

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

// TableName 表名
func (RequestLog) TableName() string {
	return "request_logs"
}
