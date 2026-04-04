package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
	"ai-gateway/internal/infras/errs"
	"ai-gateway/internal/infras/utils"
)

type APIKeyService struct {
	encryptKey  string
	apikeyRepo  repo.APIKeyRepository
	redisClient redis.UniversalClient
}

// NewAPIKeyService 创建apikey服务实例
func NewAPIKeyService(apikeyRepo repo.APIKeyRepository, redisClient redis.UniversalClient,
	encryptKey string) *APIKeyService {
	return &APIKeyService{apikeyRepo: apikeyRepo, encryptKey: encryptKey, redisClient: redisClient}
}

// Create 创建虚拟apikey信息
func (s *APIKeyService) Create(entry entity.APIKeyCreateEntity) (*entity.APIKey, string, error) {
	if entry.Prefix == "" {
		entry.Prefix = "sk"
	}

	uid := uuid.New().String()
	virtualKey := entry.Prefix + "-" + uid
	keyHash := utils.HashAPIKey(virtualKey)
	keyPrefix := utils.GetAPIKeyPrefix(virtualKey)
	key := &entity.APIKey{
		TenantID:          entry.TenantID,
		KeyHash:           keyHash,
		KeyPrefix:         keyPrefix,
		Name:              entry.Name,
		ProviderID:        entry.ProviderID,
		RateLimitPerSec:   entry.RateLimitPerSec,
		RateLimitPerMin:   entry.RateLimit,
		MonthlyTokenLimit: entry.MonthlyTokenLimit,
		TotalTokenQuota:   entry.TotalTokenQuota,
		Status:            1,
	}

	err := s.apikeyRepo.Create(key)
	return key, virtualKey, err
}

// Delete 删除apikey
func (s *APIKeyService) Delete(id uint) error {
	key, err := s.apikeyRepo.GetByID(id, "key_hash")
	if err != nil {
		return err
	}

	err = s.apikeyRepo.Delete(id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apikey:%s", utils.Md5(key.KeyHash))
	_ = s.redisClient.Del(context.Background(), cacheKey).Err()

	return nil
}

// List apikey列表
func (s *APIKeyService) List() ([]entity.APIKey, error) {
	keys, err := s.apikeyRepo.List()
	return keys, err
}

// GetByHash 根据key_hash获取租户apikey信息
func (s *APIKeyService) GetByHash(keyHash string) (*entity.APIKey, error) {
	key, err := s.apikeyRepo.GetByHash(keyHash)
	return key, err
}

// UpdateTokenConsume 更新tokens消费
// UPDATE `api_keys` SET `tokens_consumed`=tokens_consumed + 430 WHERE id = 3 AND `api_keys`.`deleted_at` IS NULL
func (s *APIKeyService) UpdateTokenConsume(id uint, tokens int) error {
	// todo 这里可以通过redis hash计数器+job定时任务，将增量数同步到db中，降低数据库压力
	return s.apikeyRepo.UpdateTokenConsume(id, tokens)
}

// GetAPIKey 获取apikey，如果缓存中不存在，就从数据库中获取
func (s *APIKeyService) GetAPIKey(ctx context.Context, apikey string) (*entity.APIKey, error) {
	keyHash := utils.HashAPIKey(apikey)
	// log.Printf("keyHash: %s", keyHash)
	cacheKey := fmt.Sprintf("apikey:%s", utils.Md5(keyHash))
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if len(val) > 0 {
		if val == "invalid" {
			return nil, errs.ErrApiKeyNotFound
		}

		var key entity.APIKey
		if json.Unmarshal([]byte(val), &key) == nil {
			return &key, nil
		}
	}

	// 从数据库中读取
	key, err := s.apikeyRepo.GetByHash(keyHash)
	if err != nil {
		if errors.Is(err, errs.ErrApiKeyNotFound) {
			_ = s.redisClient.Set(ctx, cacheKey, "invalid", 5*time.Minute).Err()
		}

		return nil, err
	}

	data, _ := json.Marshal(key)
	_ = s.redisClient.Set(ctx, cacheKey, data, 30*time.Minute).Err()

	return key, nil
}

// PerSecLimited 每秒限制
func (s *APIKeyService) PerSecLimited(ctx context.Context, id uint, rateLimitPerSec int) (bool, error) {
	secKey := fmt.Sprintf("rate_limit_sec:%d", id)
	secCount, err := s.redisClient.Incr(ctx, secKey).Result()
	if err != nil {
		log.Printf("failed to incr rate_limit_sec err:%s", err.Error())
		return true, errs.ErrServerInternal
	}

	if secCount == 1 {
		_ = s.redisClient.Expire(ctx, secKey, time.Second).Err()
	}

	if secCount > int64(rateLimitPerSec) {
		return true, fmt.Errorf("超过每秒%d次限制", rateLimitPerSec)
	}

	return false, nil
}

// PerMinLimited 每分钟频率限制
func (s *APIKeyService) PerMinLimited(ctx context.Context, id uint, rateLimitPerMin int) (bool, error) {
	rateKey := fmt.Sprintf("rate_limit_min:%d", id)
	count, err := s.redisClient.Incr(ctx, rateKey).Result()
	if err != nil {
		log.Printf("failed to incr rate_limit_min err:%s", err.Error())
		return true, errs.ErrServerInternal
	}
	if count == 1 {
		_ = s.redisClient.Expire(ctx, rateKey, time.Minute).Err()
	}

	if count > int64(rateLimitPerMin) {
		return true, fmt.Errorf("超过每分钟%d次限制", rateLimitPerMin)
	}

	return false, nil
}

// PerMonTokensLimited 每月tokens限制
func (s *APIKeyService) PerMonTokensLimited(ctx context.Context, id uint, monthlyTokenLimit int) (bool, error) {
	monthKey := fmt.Sprintf("monthly_tokens:%d", id)
	totalTokens, _ := s.redisClient.Get(ctx, monthKey).Int64()
	if totalTokens >= int64(monthlyTokenLimit) {
		return true, fmt.Errorf("已超过每月%d Token限额", monthlyTokenLimit)
	}

	return false, nil
}
