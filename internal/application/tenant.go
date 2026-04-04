package application

import (
	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
)

type TenantService struct {
	tenantRepo repo.TenantRepository
}

// NewTenantService 创建repo service
func NewTenantService(tenantRepo repo.TenantRepository) *TenantService {
	return &TenantService{tenantRepo: tenantRepo}
}

// Create 创建租户
func (s *TenantService) Create(name string) error {
	return s.tenantRepo.Create(name)
}

// List 租户列表
func (s *TenantService) List() ([]entity.Tenant, error) {
	tenants, err := s.tenantRepo.List()
	return tenants, err
}

// Delete 删除租户
func (s *TenantService) Delete(id uint) error {
	return s.tenantRepo.Delete(id)
}

// UpdateStatus 更新租户状态
func (s *TenantService) UpdateStatus(id uint, status int) error {
	return s.tenantRepo.UpdateStatus(id, status)
}
