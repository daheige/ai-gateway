package service

import (
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
)

type TenantService struct {
	db *gorm.DB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{db: db}
}

// Create 创建租户
func (s *TenantService) Create(name string) error {
	tenant := entity.Tenant{Name: name, Status: 1}
	return s.db.Create(&tenant).Error
}

// List 租户列表
func (s *TenantService) List() ([]entity.Tenant, error) {
	var tenants []entity.Tenant
	err := s.db.Find(&tenants).Error
	return tenants, err
}

// Delete 删除租户
func (s *TenantService) Delete(id uint) error {
	return s.db.Delete(&entity.Tenant{}, id).Error
}

// UpdateStatus 更新租户状态
func (s *TenantService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&entity.Tenant{}).Where("id = ?", id).Update("status", status).Error
}
