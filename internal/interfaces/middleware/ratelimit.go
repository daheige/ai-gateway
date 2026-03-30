package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/infras/utils"
)

func getAPIKeyFromCache(ctx context.Context, rdb redis.UniversalClient, db *gorm.DB,
	keyHash string) (*entity.APIKey, error) {
	cacheKey := fmt.Sprintf("apikey:%s", keyHash)
	val, err := rdb.Get(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(val) > 0 {
		var key entity.APIKey
		if json.Unmarshal([]byte(val), &key) == nil {
			return &key, nil
		}
	}

	var key entity.APIKey
	if err := db.Where("key_hash = ? AND status = 1", keyHash).First(&key).Error; err != nil {
		return nil, err
	}

	data, _ := json.Marshal(key)
	rdb.Set(ctx, cacheKey, data, 5*time.Minute)
	return &key, nil
}

// RateLimitMiddleware 频率限制中间件
func RateLimitMiddleware(rdb redis.UniversalClient, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || len(auth) <= 7 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少API Key"})
			c.Abort()
			return
		}

		// 虚拟apikey
		apiKey := strings.TrimPrefix(auth, "Bearer ")
		ctx := context.Background()

		// 获取key_hash
		keyHash := utils.HashAPIKey(apiKey)
		key, err := getAPIKeyFromCache(ctx, rdb, db, keyHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的API Key"})
			c.Abort()
			return
		}

		// 每秒请求限制
		if key.RateLimitPerSec > 0 {
			secKey := fmt.Sprintf("rate_limit_sec:%d", key.ID)
			secCount, err := rdb.Incr(ctx, secKey).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "限流检查失败"})
				c.Abort()
				return
			}
			if secCount == 1 {
				rdb.Expire(ctx, secKey, time.Second)
			}
			if secCount > int64(key.RateLimitPerSec) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("超过每秒%d次限制", key.RateLimitPerSec)})
				c.Abort()
				return
			}
		}

		// 每分钟请求次数限制
		rateKey := fmt.Sprintf("rate_limit:%d", key.ID)
		count, err := rdb.Incr(ctx, rateKey).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "限流检查失败"})
			c.Abort()
			return
		}
		if count == 1 {
			rdb.Expire(ctx, rateKey, time.Minute)
		}
		if count > int64(key.RateLimit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("超过每分钟%d次限制", key.RateLimit)})
			c.Abort()
			return
		}

		// 每月Token限制（从Redis读取，由后台job写入）
		if key.MonthlyTokenLimit > 0 {
			monthKey := fmt.Sprintf("monthly_tokens:%d", key.ID)
			totalTokens, _ := rdb.Get(ctx, monthKey).Int64()
			if totalTokens >= int64(key.MonthlyTokenLimit) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("已超过每月%d Token限额", key.MonthlyTokenLimit)})
				c.Abort()
				return
			}
		}

		// 总Token额度检查
		if key.TotalTokenQuota > 0 && key.TokensConsumed >= key.TotalTokenQuota {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("已耗尽总Token额度 %d/%d", key.TokensConsumed, key.TotalTokenQuota)})
			c.Abort()
			return
		}

		c.Set("api_key", apiKey)
		c.Set("key_hash", keyHash)
		c.Next()
	}
}

// InvalidateAPIKeyCache 删除API Key缓存（创建/删除/更新时调用）
func InvalidateAPIKeyCache(ctx context.Context, rdb *redis.Client, keyHash string) {
	rdb.Del(ctx, fmt.Sprintf("apikey:%s", keyHash))
}
