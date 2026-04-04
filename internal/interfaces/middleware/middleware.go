package middleware

import (
	"github.com/gin-gonic/gin"
)

// Middlewares 中间件
type Middlewares struct {
	RateLimitMiddleware  gin.HandlerFunc
	AdminAuthMiddleware  gin.HandlerFunc
	APIKeyAuthMiddleware gin.HandlerFunc
	LogMiddleware        gin.HandlerFunc
	MetricsMiddleware    gin.HandlerFunc
}
