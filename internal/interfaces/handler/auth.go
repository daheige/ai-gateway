package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/infras/config"
	"ai-gateway/internal/infras/utils"
)

type AuthHandler struct {
	admin     config.AdminConfig
	jwtSecret string
	jwtExpire int64
}

func NewAuthHandler(admin config.AdminConfig, secret string, expire int64) *AuthHandler {
	return &AuthHandler{admin: admin, jwtSecret: secret, jwtExpire: expire}
}

// LoginRequest 管理员登陆request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 管理员登陆
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Username != h.admin.Username || req.Password != h.admin.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, _ := utils.GenerateToken(0, h.jwtSecret, 24*3600*1000000000)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
