package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
	"ai-gateway/internal/infras/utils"
)

// ProviderService 模型provider服务
type ProviderService struct {
	encryptKey   string
	providerRepo repo.ProviderRepository
	redisClient  redis.UniversalClient
}

// NewProviderService 创建provider服务
func NewProviderService(providerRepo repo.ProviderRepository, redisClient redis.UniversalClient, encryptKey string) *ProviderService {
	return &ProviderService{providerRepo: providerRepo, encryptKey: encryptKey, redisClient: redisClient}
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
func (s *ProviderService) Delete(ctx context.Context, id uint) error {
	err := s.providerRepo.Delete(id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("model_provider:%d", id)
	_ = s.redisClient.Del(ctx, key).Err()
	return err
}

// GetByID 获取provider
func (s *ProviderService) GetByID(ctx context.Context, id uint) (*entity.Provider, error) {
	key := fmt.Sprintf("model_provider:%d", id)
	str, err := s.redisClient.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if str != "" {
		provider := &entity.Provider{}
		err = json.Unmarshal([]byte(str), provider)
		if err != nil {
			return nil, err
		}

		return provider, nil
	}

	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(provider)
	_ = s.redisClient.Set(ctx, key, string(b), 3600*time.Second).Err()

	return provider, nil
}

// DecryptAPIKey 解密得到真是的provider apikey
func (s *ProviderService) DecryptAPIKey(encKey string) (string, error) {
	return utils.DecryptAPIKey(encKey, s.encryptKey)
}
