package entity

import (
	"time"

	"gorm.io/gorm"
)

// Provider 模型provider
type Provider struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Type      string         `gorm:"size:20;not null" json:"type"`
	BaseURL   string         `gorm:"size:255" json:"base_url"`
	APIKeyEnc string         `gorm:"type:text" json:"-"`      // provider apikey加密存储
	Models    string         `gorm:"type:text" json:"models"` // 支持的模型列表
	Status    int            `gorm:"default:1" json:"status"` // 状态,1可用，0禁用
	Priority  int            `gorm:"default:0" json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 逻辑删除
}
