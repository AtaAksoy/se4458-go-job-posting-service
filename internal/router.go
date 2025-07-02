package internal

import (
	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/v1/jobs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(jobHandler *jobs.JobHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		jobsGroup := api.Group("/jobs")
		{
			jobsGroup.POST("", jobHandler.CreateJob)
			jobsGroup.GET("", jobHandler.ListJobs)
			jobsGroup.GET(":id", jobHandler.GetJobByID)
			jobsGroup.PUT(":id", jobHandler.UpdateJob)
			jobsGroup.DELETE(":id", jobHandler.DeleteJob)
			jobsGroup.GET("/search", jobHandler.SearchJobs)
		}
	}

	return r
}
