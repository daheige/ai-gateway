package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"ai-gateway/internal/application"
	"ai-gateway/internal/infras/config"
	"ai-gateway/internal/interfaces/handler"
	"ai-gateway/internal/interfaces/middleware"
	"ai-gateway/internal/interfaces/routers"
	"ai-gateway/internal/providers"
)

func main() {
	cfg := config.Load()

	// 数据库连接
	db := config.InitDB(cfg.Database)
	defer config.CloseDB(db)

	rdb := config.InitRedis(cfg.Redis)
	defer func() {
		_ = rdb.Close()
	}()

	repos := providers.NewRepositories(db)
	// 创建服务
	apiKeyService := application.NewAPIKeyService(repos.APIKeyRepo, rdb, cfg.Encrypt.Key)
	providerService := application.NewProviderService(repos.ProviderRepo, rdb, cfg.Encrypt.Key)
	tenantService := application.NewTenantService(repos.TenantRepo)
	statsService := application.NewStatsService(repos.StatsRepo)
	logService := application.NewLogService(repos.LogRepo)

	// 初始化handlers
	handlers := handler.Handlers{
		AuthHandler:     handler.NewAuthHandler(cfg.Admin, cfg.JWT.Secret, int64(cfg.JWT.ExpireTime)),
		ApiKeyHandler:   handler.NewAPIKeyHandler(apiKeyService),
		GatewayHandler:  handler.NewGatewayHandler(apiKeyService, providerService, logService),
		LogHandler:      handler.NewLogHandler(logService),
		ProviderHandler: handler.NewProviderHandler(providerService),
		TenantHandler:   handler.NewTenantHandler(tenantService),
		MetricsHandler:  handler.NewMetricsHandler(),
		StatsHandler:    handler.NewStatsHandler(statsService),
	}

	// 初始化中间件
	middlewares := middleware.Middlewares{
		APIKeyAuthMiddleware: middleware.APIKeyAuthMiddleware(apiKeyService),
		RateLimitMiddleware:  middleware.RateLimitMiddleware(apiKeyService),
		AdminAuthMiddleware:  middleware.AdminAuthMiddleware(cfg.JWT.Secret),
		LogMiddleware:        middleware.LogMiddleware(),
		MetricsMiddleware:    middleware.MetricsMiddleware(),
	}

	// 启动http server
	r := gin.New()
	routers.NewRouter(r, handlers, middlewares)

	// 启动服务
	address := fmt.Sprintf("0.0.0.0:%d", cfg.AppPort)
	server := &http.Server{
		Handler:           r,
		Addr:              address,
		ReadHeaderTimeout: 10 * time.Second, // read header timeout
		ReadTimeout:       10 * time.Second, // read request timeout
		WriteTimeout:      30 * time.Second, // write timeout
		IdleTimeout:       20 * time.Second, // tcp idle time
	}

	// 在独立携程中运行
	log.Printf("server listening on %s\n", address)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("server close error", map[string]interface{}{
					"trace_error": err.Error(),
				})
				return
			}

			log.Println("server will exit...")
		}
	}()

	// 等待平滑退出
	shutdown(server, cfg.GracefulWait)
}

func shutdown(server *http.Server, gracefulWait time.Duration) {
	// server平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// receive signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// linux signal,please use this in production.
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), gracefulWait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		err := server.Shutdown(ctx)
		if err != nil {
			log.Println("server shutdown error:", err)
		}
	}()

	select {
	case <-done:
		log.Println("server shutting down")
	case <-ctx.Done():
		log.Println("server ctx timeout")
	}
}
