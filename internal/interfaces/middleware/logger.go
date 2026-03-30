package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LogMiddleware 日志中间件
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("[Gateway] %s %s %d %v %s",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency, c.ClientIP())
	}
}
