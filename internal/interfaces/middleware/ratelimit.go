package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
	"ai-gateway/internal/domain/entity"
	"ai-gateway/internal/infras/ctxkeys"
)

// RateLimitMiddleware 频率限制中间件
func RateLimitMiddleware(apiKeySvc *application.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		key, _ := ctx.Value(ctxkeys.APIKeyInfo).(*entity.APIKey)
		// 每秒请求限制
		if key.RateLimitPerSec > 0 {
			limited, err := apiKeySvc.PerSecLimited(ctx, key.ID, key.RateLimitPerSec)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if limited {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": fmt.Sprintf("超过每秒%d次限制", key.RateLimitPerSec),
				})
				return
			}
		}

		// 每分钟请求次数限制
		if key.RateLimitPerMin > 0 {
			limited, err := apiKeySvc.PerMinLimited(ctx, key.ID, key.RateLimitPerMin)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if limited {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": fmt.Sprintf("超过每分%d次限制", key.RateLimitPerMin),
				})
				return
			}
		}

		// 每月Token限制（从Redis读取，由后台job写入）
		if key.MonthlyTokenLimit > 0 {
			limited, err := apiKeySvc.PerMonTokensLimited(ctx, key.ID, key.MonthlyTokenLimit)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if limited {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": fmt.Sprintf("已超过每月%d Token限额", key.MonthlyTokenLimit),
				})
				return
			}
		}

		// 总Token额度检查
		if key.TotalTokenQuota > 0 && key.TokensConsumed >= key.TotalTokenQuota {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("已耗尽总Token额度 %d/%d", key.TokensConsumed, key.TotalTokenQuota),
			})
			return
		}

		c.Next()
	}
}
