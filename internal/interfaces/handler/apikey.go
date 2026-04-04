package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
	"ai-gateway/internal/domain/entity"
)

// APIKeyHandler apikey handler
type APIKeyHandler struct {
	service *application.APIKeyService
}

// NewAPIKeyHandler 创建apikey handler
func NewAPIKeyHandler(service *application.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: service}
}

// APIKeyCreateRequest 虚拟apikey创建请求
type APIKeyCreateRequest struct {
	Name              string `json:"name" binding:"required"`
	TenantID          uint   `json:"tenant_id" binding:"required"`
	ProviderID        uint   `json:"provider_id" binding:"required"`
	Prefix            string `json:"prefix"`
	RateLimitPerSec   int    `json:"rate_limit_per_sec"`
	RateLimit         int    `json:"rate_limit"`
	MonthlyTokenLimit int    `json:"monthly_token_limit"`
	TotalTokenQuota   int    `json:"total_token_quota"`
}

// Create 创建虚拟apikey
func (h *APIKeyHandler) Create(c *gin.Context) {
	var req APIKeyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry := entity.APIKeyCreateEntity{
		TenantID:          req.TenantID,
		ProviderID:        req.ProviderID,
		Name:              req.Name,
		Prefix:            req.Prefix,
		RateLimitPerSec:   req.RateLimitPerSec,
		RateLimit:         req.RateLimit,
		MonthlyTokenLimit: req.MonthlyTokenLimit,
		TotalTokenQuota:   req.TotalTokenQuota,
	}
	key, virtualKey, err := h.service.Create(entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": key.ID, "key": virtualKey, "key_prefix": key.KeyPrefix})
}

// Delete 删除虚拟apikey
func (h *APIKeyHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// List apikey 列表
func (h *APIKeyHandler) List(c *gin.Context) {
	keys, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, keys)
}
