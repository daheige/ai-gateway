package routers

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	gateway "ai-gateway/internal/interfaces/embed"
	"ai-gateway/internal/interfaces/handler"
	"ai-gateway/internal/interfaces/middleware"
)

// NewRouter 创建路由
func NewRouter(r *gin.Engine, handlers handler.Handlers, m middleware.Middlewares) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, ai-gateway")
	})

	// 嵌入静态文件
	staticSub, _ := fs.Sub(gateway.StaticFS, "static")
	r.GET("/admin/*filepath", gin.WrapH(http.StripPrefix("/admin", http.FileServer(http.FS(staticSub)))))

	r.GET("/metrics", handlers.MetricsHandler.GetMetrics)

	r.POST("/auth/login", handlers.AuthHandler.Login)

	api := r.Group("/api")
	api.Use(m.AuthMiddleware)
	{
		api.POST("/keys", handlers.ApiKeyHandler.Create)
		api.DELETE("/keys/:id", handlers.ApiKeyHandler.Delete)
		api.GET("/keys", handlers.ApiKeyHandler.List)
		api.GET("/logs", handlers.LogHandler.List)
		api.POST("/providers", handlers.ProviderHandler.Create)
		api.GET("/providers", handlers.ProviderHandler.List)
		api.DELETE("/providers/:id", handlers.ProviderHandler.Delete)
		api.POST("/tenants", handlers.TenantHandler.Create)
		api.GET("/tenants", handlers.TenantHandler.List)
		api.PUT("/tenants/:id", handlers.TenantHandler.Update)
		api.DELETE("/tenants/:id", handlers.TenantHandler.Delete)
		api.GET("/stats/overview", handlers.StatsHandler.Overview)
		api.GET("/stats/daily", handlers.StatsHandler.DailyStats)
		api.GET("/stats/keys", handlers.StatsHandler.KeyStats)
	}

	gw := r.Group("/v1")
	gw.Use(m.MetricsMiddleware)
	gw.Use(m.RateLimitMiddleware)
	gw.Use(m.LogMiddleware)
	{
		gw.POST("/chat/completions", handlers.GatewayHandler.Proxy)
	}
}
