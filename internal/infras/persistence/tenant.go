package persistence

import (
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

// TenantRepoImpl 租户接口实现
type TenantRepoImpl struct {
	db *gorm.DB
}

var _ repo.TenantRepository = (*TenantRepoImpl)(nil)

// NewTenantRepo 创建租户repo接口
func NewTenantRepo(db *gorm.DB) repo.TenantRepository {
	return &TenantRepoImpl{db: db}
}

// Create 创建租户
func (s *TenantRepoImpl) Create(name string) error {
	tenant := entity.Tenant{Name: name, Status: 1}
	return s.db.Create(&tenant).Error
}

// List 租户列表
func (s *TenantRepoImpl) List() ([]entity.Tenant, error) {
	var tenants []entity.Tenant
	err := s.db.Find(&tenants).Error
	return tenants, err
}

// Delete 删除租户
func (s *TenantRepoImpl) Delete(id uint) error {
	return s.db.Delete(&entity.Tenant{}, id).Error
}

// UpdateStatus 更新租户状态
func (s *TenantRepoImpl) UpdateStatus(id uint, status int) error {
	return s.db.Model(&entity.Tenant{}).Where("id = ?", id).Update("status", status).Error
}
