package persistence

import (
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

// ProviderRepoImpl 模型provider实现
type ProviderRepoImpl struct {
	db *gorm.DB
}

var _ repo.ProviderRepository = (*ProviderRepoImpl)(nil)

// NewProviderRepo 初始化provider repo实例
func NewProviderRepo(db *gorm.DB) repo.ProviderRepository {
	return &ProviderRepoImpl{
		db: db,
	}
}

// Create 创建provider
func (p *ProviderRepoImpl) Create(provider *entity.Provider) error {
	return p.db.Create(&provider).Error
}

// List provider 列表
func (p *ProviderRepoImpl) List() ([]entity.Provider, error) {
	var providers []entity.Provider
	err := p.db.Find(&providers).Error
	return providers, err
}

// Delete 删除provider
func (p *ProviderRepoImpl) Delete(id uint) error {
	err := p.db.Where("id = ?", id).Delete(&entity.Provider{}).Error
	return err
}

// GetByID 根据id获取provider
func (p *ProviderRepoImpl) GetByID(id uint) (*entity.Provider, error) {
	provider := &entity.Provider{}
	err := p.db.Where("id = ? AND status = ?", id, 1).First(&provider).Error
	return provider, err
}
