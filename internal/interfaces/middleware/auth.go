package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
	"ai-gateway/internal/infras/ctxkeys"
	"ai-gateway/internal/infras/utils"
)

// AdminAuthMiddleware jwt授权中间件
// 用于管理后台请求认证
func AdminAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := utils.ParseToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
			c.Abort()
			return
		}

		c.Set("tenant_id", claims.TenantID)
		c.Next()
	}
}

// APIKeyAuthMiddleware apikey认证中间件
func APIKeyAuthMiddleware(apiKeySvc *application.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || len(auth) <= 7 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少API Key"})
			return
		}

		// 虚拟apikey
		apiKey := strings.TrimPrefix(auth, "Bearer ")
		key, err := apiKeySvc.GetAPIKey(utils.NewContext(ctx), apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的API Key"})
			return
		}
		
		requestID, _ := ctx.Value(ctxkeys.XRequestID).(string)
		log.Printf(
			"apikey auth success,request_id:%s tenant_id:%d apikey:%s",
			requestID, key.TenantID, utils.GetAPIKeyPrefix(apiKey),
		)

		// 将必要的信息放入上下文中
		ctx = context.WithValue(ctx, ctxkeys.APIKeyInfo, key)
		ctx = context.WithValue(ctx, ctxkeys.TenantID, key.TenantID)
		ctx = context.WithValue(ctx, ctxkeys.APIKey, apiKey)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
