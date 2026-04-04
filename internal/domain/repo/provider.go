package repo

import (
	"ai-gateway/internal/domain/entity"
)

// ProviderRepository 模型provider接口
type ProviderRepository interface {
	Create(provider *entity.Provider) error
	List() ([]entity.Provider, error)
	Delete(id uint) error
	GetByID(id uint) (*entity.Provider, error)
}
