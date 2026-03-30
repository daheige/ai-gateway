package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
)

type LogHandler struct {
	db *gorm.DB
}

func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{db: db}
}

type LogAuditRequest struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

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

	var total int64
	h.db.Model(&entity.RequestLog{}).Count(&total)

	var logs []entity.RequestLog
	h.db.Order("created_at desc").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  req.Page,
		"limit": req.Limit,
		"data":  logs,
	})
}
