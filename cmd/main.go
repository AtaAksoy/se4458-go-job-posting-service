// @title           Job Posting API
// @version         1.0
// @description     Job posting service with CRUD and pagination.
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"context"
	"log"

	"github.com/AtaAksoy/se4458-go-job-posting-service/config"
	_ "github.com/AtaAksoy/se4458-go-job-posting-service/docs"
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal"
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/v1/db"
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/v1/jobs"
)

func main() {
	cfg := config.LoadConfig()

	dbConn := db.Connect(cfg.DBDSN, &jobs.Job{})

	redisClient := db.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)

	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	jobCache := jobs.NewJobCache(redisClient)

	repo := jobs.NewGormJobRepository(dbConn, jobCache)
	handler := jobs.NewJobHandler(repo)

	r := internal.SetupRouter(handler)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
