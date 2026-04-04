package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/infras/ctxkeys"
	"ai-gateway/internal/infras/utils"
)

// LogMiddleware 日志中间件
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = utils.UUID()
		}

		ctx := c.Request.Context()
		ip := c.ClientIP()
		uri := c.Request.RequestURI
		method := c.Request.Method
		ctx = context.WithValue(ctx, ctxkeys.XRequestID, requestID)
		ctx = context.WithValue(ctx, ctxkeys.RequestURI, uri)
		ctx = context.WithValue(ctx, ctxkeys.RequestMethod, c.Request.Method)
		ctx = context.WithValue(ctx, ctxkeys.ClientIP, ip)
		ctx = context.WithValue(ctx, ctxkeys.UserAgent, c.Request.UserAgent())
		c.Request = c.Request.WithContext(ctx)

		// 记录开始请求日志
		log.Printf("[Gateway] exec begin,request_id:%s method:%s uri:%s ip:%s", requestID, method, uri, ip)
		c.Next()

		// 记录结束请求日志
		c.Header("X-Request-Id", requestID)
		latency := time.Since(start)
		log.Printf(
			"[Gateway] exec end,request_id:%s method:%s uri:%s http_status:%d latency:%v",
			requestID, method, uri, c.Writer.Status(), latency,
		)
	}
}
