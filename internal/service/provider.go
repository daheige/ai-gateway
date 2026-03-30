package service

import (
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/infras/utils"
)

type ProviderService struct {
	db         *gorm.DB
	encryptKey string
}

func NewProviderService(db *gorm.DB, encryptKey string) *ProviderService {
	return &ProviderService{db: db, encryptKey: encryptKey}
}

// Create 创建模型提供商provider
func (s *ProviderService) Create(name, providerType, baseURL, apiKey, modelList string, priority int) error {
	encAPIKey, err := utils.EncryptAPIKey(apiKey, s.encryptKey)
	if err != nil {
		return err
	}
	provider := entity.Provider{
		Name:      name,
		Type:      providerType,
		BaseURL:   baseURL,
		APIKeyEnc: encAPIKey,
		Models:    modelList,
		Priority:  priority,
		Status:    1,
	}
	return s.db.Create(&provider).Error
}

// List 列出providers
func (s *ProviderService) List() ([]entity.Provider, error) {
	var providers []entity.Provider
	err := s.db.Find(&providers).Error
	return providers, err
}

// Delete 删除provider
func (s *ProviderService) Delete(id uint) error {
	return s.db.Where("id = ?", id).Delete(&entity.Provider{}).Error
}

// GetByID 获取provider
func (s *ProviderService) GetByID(id uint) (*entity.Provider, error) {
	var provider entity.Provider
	err := s.db.Where("id = ? AND status = 1", id).First(&provider).Error
	return &provider, err
}

// DecryptAPIKey 解密得到真是的provider apikey
func (s *ProviderService) DecryptAPIKey(encKey string) (string, error) {
	return utils.DecryptAPIKey(encKey, s.encryptKey)
}
