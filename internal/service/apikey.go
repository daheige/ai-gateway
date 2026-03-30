package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/infras/utils"
)

type APIKeyService struct {
	db         *gorm.DB
	encryptKey string
}

// NewAPIKeyService 创建apikey服务实例
func NewAPIKeyService(db *gorm.DB, encryptKey string) *APIKeyService {
	return &APIKeyService{db: db, encryptKey: encryptKey}
}

// Create 创建虚拟apikey信息
func (s *APIKeyService) Create(entry entity.APIKeyCreateEntity) (*entity.APIKey, string, error) {
	if entry.Prefix == "" {
		entry.Prefix = "sk"
	}

	uid := uuid.New().String()
	uid = uid[:8] + uid[9:13] + uid[14:18] + uid[19:23] + uid[24:]
	virtualKey := entry.Prefix + "-" + uid
	keyHash := utils.HashAPIKey(virtualKey)
	keyPrefix := utils.GetAPIKeyPrefix(virtualKey)
	key := entity.APIKey{
		TenantID:          entry.TenantID,
		KeyHash:           keyHash,
		KeyPrefix:         keyPrefix,
		Name:              entry.Name,
		ProviderID:        entry.ProviderID,
		RateLimitPerSec:   entry.RateLimitPerSec,
		RateLimit:         entry.RateLimit,
		MonthlyTokenLimit: entry.MonthlyTokenLimit,
		TotalTokenQuota:   entry.TotalTokenQuota,
		Status:            1,
	}

	err := s.db.Create(&key).Error
	return &key, virtualKey, err
}

func (s *APIKeyService) Delete(id uint) error {
	return s.db.Where("id = ?", id).Delete(&entity.APIKey{}).Error
}

func (s *APIKeyService) List() ([]entity.APIKey, error) {
	var keys []entity.APIKey
	err := s.db.Find(&keys).Error
	return keys, err
}

// GetByHash 根据key_hash获取租户apikey信息
func (s *APIKeyService) GetByHash(keyHash string) (*entity.APIKey, error) {
	var key entity.APIKey
	err := s.db.Where("key_hash = ? AND status = 1", keyHash).First(&key).Error
	return &key, err
}

// DecryptProviderKey 解密apikey
func (s *APIKeyService) DecryptProviderKey(encKey string) (string, error) {
	return utils.DecryptAPIKey(encKey, s.encryptKey)
}
