package persistence

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/domain/repo"
	"ai-gateway/internal/infras/errs"
)

// APIKeyRepoImpl 虚拟apikey接口实现
type APIKeyRepoImpl struct {
	db *gorm.DB
}

var _ repo.APIKeyRepository = (*APIKeyRepoImpl)(nil)

// NewAPIKeyRepo 创建apikey接口实例
func NewAPIKeyRepo(db *gorm.DB) repo.APIKeyRepository {
	return &APIKeyRepoImpl{db: db}
}

// Create 创建虚拟apikey信息
func (s *APIKeyRepoImpl) Create(entry *entity.APIKey) error {
	err := s.db.Create(&entry).Error
	return err
}

func (s *APIKeyRepoImpl) Delete(id uint) error {
	return s.db.Where("id = ?", id).Delete(&entity.APIKey{}).Error
}

func (s *APIKeyRepoImpl) List() ([]entity.APIKey, error) {
	var keys []entity.APIKey
	err := s.db.Find(&keys).Error
	return keys, err
}

// GetByHash 根据key_hash获取租户apikey信息
func (s *APIKeyRepoImpl) GetByHash(keyHash string) (*entity.APIKey, error) {
	key := &entity.APIKey{}
	err := s.db.Where("key_hash = ? AND status = 1", keyHash).First(key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrApiKeyNotFound
		}

		return nil, err
	}

	return key, nil
}

// UpdateTokenConsume 更新token消费
func (s *APIKeyRepoImpl) UpdateTokenConsume(id uint, tokens int) error {
	err := s.db.Model(&entity.APIKey{}).Where("id = ?", id).
		UpdateColumn("tokens_consumed", gorm.Expr("tokens_consumed + ?", tokens)).Error
	return err
}

// GetByID 通过id获取apikey
func (s *APIKeyRepoImpl) GetByID(id uint, cols ...string) (*entity.APIKey, error) {
	if len(cols) == 0 {
		cols = []string{"*"}
	}

	key := &entity.APIKey{}
	err := s.db.Where("id = ?", id).Select(strings.Join(cols, ",")).First(key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrApiKeyNotFound
		}

		return nil, err
	}

	return key, nil
}
