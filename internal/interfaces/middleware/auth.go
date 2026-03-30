package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/infras/utils"
)

// AuthMiddleware jwt授权中间件
// 用于管理后台请求认证
func AuthMiddleware(secret string) gin.HandlerFunc {
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
