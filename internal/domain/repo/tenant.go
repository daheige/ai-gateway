package repo

import (
	"ai-gateway/internal/domain/entity"
)

// TenantRepository 租户接口定义
type TenantRepository interface {
	Create(username string) error
	List() ([]entity.Tenant, error)
	Delete(id uint) error
	UpdateStatus(id uint, status int) error
}
