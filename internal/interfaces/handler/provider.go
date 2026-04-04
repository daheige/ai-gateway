package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
)

type ProviderHandler struct {
	service *application.ProviderService
}

// NewProviderHandler 创建provider handler
func NewProviderHandler(service *application.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

type ProviderRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Models   string `json:"models"`
	Priority int    `json:"priority"`
}

// Create 创建provider
func (h *ProviderHandler) Create(c *gin.Context) {
	var req ProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(req.Name, req.Type, req.BaseURL, req.APIKey, req.Models, req.Priority); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// List provider 列表
func (h *ProviderHandler) List(c *gin.Context) {
	providers, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, providers)
}

// Delete 删除provider
func (h *ProviderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()
	if err := h.service.Delete(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
