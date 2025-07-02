package jobs

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	repo JobRepository
}

func NewJobHandler(repo JobRepository) *JobHandler {
	return &JobHandler{repo: repo}
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job := Job{
		Title:       req.Title,
		Description: req.Description,
		Company:     req.Company,
		City:        req.City,
		State:       req.State,
		Status:      true,
		CreatedAt:   time.Now().Unix(),
	}
	if err := h.repo.Create(&job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}
	c.JSON(http.StatusCreated, JobResponse{
		ID:          job.ID,
		Title:       job.Title,
		Description: job.Description,
		Company:     job.Company,
		City:        job.City,
		State:       job.State,
		CreatedAt:   job.CreatedAt,
		Status:      job.Status,
	})
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	jobs, total, err := h.repo.List(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
		return
	}
	responses := make([]JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = JobResponse{
			ID:          job.ID,
			Title:       job.Title,
			Description: job.Description,
			Company:     job.Company,
			City:        job.City,
			State:       job.State,
			CreatedAt:   job.CreatedAt,
			Status:      job.Status,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"jobs":  responses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job id"})
		return
	}
	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}
	c.Status(http.StatusNoContent)
}
