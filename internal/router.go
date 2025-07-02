package internal

import (
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/jobs"

	"github.com/gin-gonic/gin"
)

func SetupRouter(jobHandler *jobs.JobHandler) *gin.Engine {
	r := gin.Default()

	jobsGroup := r.Group("/jobs")
	{
		jobsGroup.POST("", jobHandler.CreateJob)
		jobsGroup.GET("", jobHandler.ListJobs)
		jobsGroup.DELETE(":id", jobHandler.DeleteJob)
	}

	return r
}
