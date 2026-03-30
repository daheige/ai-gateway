package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/interfaces/middleware"
)

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metrics := middleware.GetMetrics()
	c.JSON(http.StatusOK, metrics)
}
