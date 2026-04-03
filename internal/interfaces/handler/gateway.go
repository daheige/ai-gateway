package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ai-gateway/internal/domain/entity"
	service2 "ai-gateway/internal/service"
)

// GatewayHandler gateway网关handler
type GatewayHandler struct {
	db              *gorm.DB
	apiKeyService   *service2.APIKeyService
	providerService *service2.ProviderService
}

// NewGatewayHandler 创建gateway handler
func NewGatewayHandler(db *gorm.DB, apiKeyService *service2.APIKeyService, providerService *service2.ProviderService) *GatewayHandler {
	return &GatewayHandler{
		db:              db,
		apiKeyService:   apiKeyService,
		providerService: providerService,
	}
}

// Proxy 转发请求
func (h *GatewayHandler) Proxy(c *gin.Context) {
	startTime := time.Now()
	keyHash := c.GetString("key_hash")
	keyEntity, err := h.apiKeyService.GetByHash(keyHash)
	if err != nil {
		log.Printf("get key by hash error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
		return
	}

	bodyBytes, _ := io.ReadAll(c.Request.Body)

	// 这里可以定义结构体进一步优化这里 todo
	var reqBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &reqBody)
	if err != nil {
		log.Printf("unmarshal request body error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	}

	model := reqBody["model"].(string)
	reqBody = nil

	provider, err := h.providerService.GetByID(keyEntity.ProviderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not found"})
		return
	}

	// 得到真正的apikey
	realAPIKey, err := h.providerService.DecryptAPIKey(provider.APIKeyEnc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decrypt failed"})
		return
	}

	url := provider.BaseURL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "chat/completions"

	// log.Println("url:", url)
	// log.Println("body:", string(bodyBytes))
	// log.Println("realAPIKey: ", realAPIKey)
	ctx := c.Request.Context()
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+realAPIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 读取body内容
	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokens, promptTokens, compTokens := 0, 0, 0
	if usage, ok := result["usage"].(map[string]interface{}); ok {
		if t, ok := usage["total_tokens"].(float64); ok {
			tokens = int(t)
		}

		if p, ok := usage["prompt_tokens"].(float64); ok {
			promptTokens = int(p)
		}
		if co, ok := usage["completion_tokens"].(float64); ok {
			compTokens = int(co)
		}
	}

	latency := int(time.Since(startTime).Milliseconds())

	// 记录请求日志
	h.db.Create(&entity.RequestLog{
		TenantID:     keyEntity.TenantID,
		APIKeyID:     keyEntity.ID,
		ProviderID:   provider.ID,
		Model:        model,
		TokensUsed:   tokens,
		PromptTokens: promptTokens,
		CompTokens:   compTokens,
		Status:       resp.StatusCode,
		Latency:      latency,
		IP:           c.ClientIP(),
	})

	if tokens > 0 {
		h.db.Model(&entity.APIKey{}).Where("id = ?", keyEntity.ID).
			UpdateColumn("tokens_consumed", gorm.Expr("tokens_consumed + ?", tokens))
	}

	c.Data(resp.StatusCode, "application/json", respBody)
}
