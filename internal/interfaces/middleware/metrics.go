package middleware

import (
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	totalRequests int64
	totalTokens   int64
	totalLatency  int64
	requestCount  int64
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start).Milliseconds()
		atomic.AddInt64(&totalRequests, 1)
		atomic.AddInt64(&totalLatency, latency)
		atomic.AddInt64(&requestCount, 1)
	}
}

func RecordTokens(tokens int) {
	atomic.AddInt64(&totalTokens, int64(tokens))
}

func GetMetrics() map[string]interface{} {
	reqs := atomic.LoadInt64(&totalRequests)
	tokens := atomic.LoadInt64(&totalTokens)
	latency := atomic.LoadInt64(&totalLatency)
	count := atomic.LoadInt64(&requestCount)

	avgLatency := int64(0)
	if count > 0 {
		avgLatency = latency / count
	}

	return map[string]interface{}{
		"total_requests": reqs,
		"total_tokens":   tokens,
		"avg_latency_ms": avgLatency,
	}
}
