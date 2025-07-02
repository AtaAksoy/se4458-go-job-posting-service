package main

import (
	"log"

	"github.com/AtaAksoy/se4458-go-job-posting-service/config"
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal"
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/db"
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/jobs"
)

func main() {
	cfg := config.LoadConfig()
	dbConn := db.Connect(cfg.DBDSN, &jobs.Job{})
	repo := jobs.NewGormJobRepository(dbConn)
	handler := jobs.NewJobHandler(repo)

	r := internal.SetupRouter(handler)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
