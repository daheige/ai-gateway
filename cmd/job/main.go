package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ai-gateway/cmd/job/tasks"
	"ai-gateway/internal/infras/config"
)

func main() {
	cfg := config.Load()

	// 初始化db
	db := config.InitDB(cfg.Database)
	defer config.CloseDB(db)

	// 初始化redis
	rdb := config.InitRedis(cfg.Redis)
	defer func() {
		_ = rdb.Close()
	}()

	job := tasks.NewTokenSyncJob(db, rdb)
	job.Start()
	log.Println("token-sync job is running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(),
		cfg.GracefulWait,
	)
	defer cancel()

	go func() {
		job.Stop()
	}()

	select {
	case <-ctx.Done():
		log.Printf("job shutdown ctx cancel error: %v", ctx.Err())
	default:
		log.Printf("job shutdown success")
	}
}
