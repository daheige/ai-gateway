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
	APIKeyEnc string         `gorm:"type:text" json:"-"`
	Models    string         `gorm:"type:text" json:"models"`
	Status    int            `gorm:"default:1" json:"status"`
	Priority  int            `gorm:"default:0" json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
