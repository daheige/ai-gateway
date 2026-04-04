package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/interfaces/middleware"
)

// MetricsHandler metrics 处理器
type MetricsHandler struct{}

// NewMetricsHandler 创建 metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// GetMetrics 获取指标
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metrics := middleware.GetMetrics()
	c.JSON(http.StatusOK, metrics)
}
