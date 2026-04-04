package application

import (
	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
	"ai-gateway/internal/infras/utils"
)

type ProviderService struct {
	encryptKey   string
	providerRepo repo.ProviderRepository
}

func NewProviderService(providerRepo repo.ProviderRepository, encryptKey string) *ProviderService {
	return &ProviderService{providerRepo: providerRepo, encryptKey: encryptKey}
}

// Create 创建模型提供商provider
func (s *ProviderService) Create(name, providerType, baseURL, apiKey, modelList string, priority int) error {
	// 加密apikey
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

	err = s.providerRepo.Create(&provider)
	return err
}

// List 列出providers
func (s *ProviderService) List() ([]entity.Provider, error) {
	providers, err := s.providerRepo.List()
	return providers, err
}

// Delete 删除provider
func (s *ProviderService) Delete(id uint) error {
	err := s.providerRepo.Delete(id)
	return err
}

// GetByID 获取provider
func (s *ProviderService) GetByID(id uint) (*entity.Provider, error) {
	provider, err := s.providerRepo.GetByID(id)
	return provider, err
}

// DecryptAPIKey 解密得到真是的provider apikey
func (s *ProviderService) DecryptAPIKey(encKey string) (string, error) {
	return utils.DecryptAPIKey(encKey, s.encryptKey)
}
