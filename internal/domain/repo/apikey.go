package repo

import (
	"ai-gateway/internal/domain/entity"
)

// APIKeyRepository 虚拟apikey接口
type APIKeyRepository interface {
	Create(entry entity.APIKeyCreateEntity) (*entity.APIKey, string, error)
	Delete(id uint) error
	List() ([]entity.APIKey, error)
	GetByHash(keyHash string) (*entity.APIKey, error)
}
