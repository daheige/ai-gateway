package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
)

// LogHandler 创建日志handler
type LogHandler struct {
	logService *application.LogService
}

// NewLogHandler 创建log handler
func NewLogHandler(logService *application.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

// LogAuditRequest 设计日志请求
type LogAuditRequest struct {
	Limit int `form:"limit" binding:"omitempty" json:"limit"`
	Page  int `form:"page" binding:"required" json:"page"`
}

// List 日志列表
func (h *LogHandler) List(c *gin.Context) {
	var req LogAuditRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	logs, total, err := h.logService.List(req.Limit, req.Page)
	if err != nil {
		// log.Println("err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  req.Page,
		"limit": req.Limit,
		"data":  logs,
	})
}
