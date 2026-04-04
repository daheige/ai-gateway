package repo

import (
	"ai-gateway/internal/domain/entity"
)

// APIKeyRepository 虚拟apikey接口
type APIKeyRepository interface {
	Create(entry *entity.APIKey) error
	Delete(id uint) error
	GetByID(id uint, cols ...string) (*entity.APIKey, error)
	List() ([]entity.APIKey, error)
	GetByHash(keyHash string) (*entity.APIKey, error)
	
	// UpdateTokenConsume 更新token消费
	UpdateTokenConsume(id uint, tokens int) error
}
