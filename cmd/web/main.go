package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"ai-gateway/internal/infras/config"
	"ai-gateway/internal/interfaces/handler"
	"ai-gateway/internal/interfaces/middleware"
	"ai-gateway/internal/routers"
	service2 "ai-gateway/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// db.AutoMigrate(&entity.Tenant{}, &entity.APIKey{}, &entity.RequestLog{}, &entity.TokenUsage{}, &entity.Provider{})
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{
			cfg.Redis.Addr,
		},
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	apiKeyService := service2.NewAPIKeyService(db, cfg.Encrypt.Key)
	providerService := service2.NewProviderService(db, cfg.Encrypt.Key)
	tenantService := service2.NewTenantService(db)

	authHandler := handler.NewAuthHandler(cfg.Admin, cfg.JWT.Secret, int64(cfg.JWT.ExpireTime))
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyService)
	gatewayHandler := handler.NewGatewayHandler(db, apiKeyService, providerService)
	logHandler := handler.NewLogHandler(db)
	providerHandler := handler.NewProviderHandler(providerService)
	tenantHandler := handler.NewTenantHandler(tenantService)
	metricsHandler := handler.NewMetricsHandler()
	statsService := service2.NewStatsService(db)
	statsHandler := handler.NewStatsHandler(statsService)

	// 初始化handlers
	handlers := handler.Handlers{
		AuthHandler:     authHandler,
		ApiKeyHandler:   apiKeyHandler,
		GatewayHandler:  gatewayHandler,
		LogHandler:      logHandler,
		ProviderHandler: providerHandler,
		TenantHandler:   tenantHandler,
		MetricsHandler:  metricsHandler,
		StatsHandler:    statsHandler,
	}

	// 初始化中间件
	middlewares := middleware.Middlewares{
		RateLimitMiddleware: middleware.RateLimitMiddleware(rdb, db),
		AuthMiddleware:      middleware.AuthMiddleware(cfg.JWT.Secret),
		LogMiddleware:       middleware.LogMiddleware(),
		MetricsMiddleware:   middleware.MetricsMiddleware(),
	}

	// 启动http server
	r := gin.New()
	r.Use(gin.Recovery())

	routers.NewRouter(r, handlers, middlewares)
	err = r.Run(cfg.Server.Port)
	if err != nil {
		log.Fatal("service run err:", err)
	}
}
