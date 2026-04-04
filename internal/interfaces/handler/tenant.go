package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
)

// TenantHandler 租户handler
type TenantHandler struct {
	service *application.TenantService
}

// NewTenantHandler 创建租户handler
func NewTenantHandler(service *application.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

// TenantRequest 租户请求结构体
type TenantRequest struct {
	Name string `json:"name" binding:"required"`
}

// Create 创建租户
func (h *TenantHandler) Create(c *gin.Context) {
	var req TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// List 租户列表
func (h *TenantHandler) List(c *gin.Context) {
	tenants, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

// Delete 删除租户
func (h *TenantHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// TenantStatusRequest 租户状态请求
type TenantStatusRequest struct {
	Status int `json:"status"`
}

// Update 更新租户状态
func (h *TenantHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req TenantStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
