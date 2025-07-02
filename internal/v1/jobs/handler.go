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

// CreateJob godoc
// @Summary      Create a new job
// @Description  Add a new job posting
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        job  body  CreateJobRequest  true  "Job info"
// @Success      201  {object}  JobResponse
// @Failure      400  {object}  map[string]string
// @Router       /jobs [post]
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
	if err := h.repo.Create(c.Request.Context(), &job); err != nil {
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

// ListJobs godoc
// @Summary      List jobs
// @Description  Get jobs with pagination
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        page   query     int false "Page number"
// @Param        limit  query     int false "Page size"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /jobs [get]
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
	jobs, total, err := h.repo.List(c.Request.Context(), offset, limit)
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

// DeleteJob godoc
// @Summary      Delete a job
// @Description  Delete a job by ID
// @Tags         jobs
// @Param        id   path      int  true  "Job ID"
// @Success      204  {string}  string  ""
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /jobs/{id} [delete]
func (h *JobHandler) DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job id"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}
	c.Status(http.StatusNoContent)
}

// SearchJobs godoc
// @Summary      Search jobs
// @Description  Search jobs by query string in title, description, company, city, or state
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        q      query     string true  "Search query"
// @Param        page   query     int    false "Page number"
// @Param        limit  query     int    false "Page size"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /jobs/search [get]
func (h *JobHandler) SearchJobs(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query"})
		return
	}
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
	jobs, total, err := h.repo.Search(c.Request.Context(), q, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search jobs"})
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

// GetJobByID godoc
// @Summary      Get a job by ID
// @Description  Get a specific job by its ID
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      200  {object}  JobResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /jobs/{id} [get]
func (h *JobHandler) GetJobByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job id"})
		return
	}
	job, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job"})
		return
	}
	c.JSON(http.StatusOK, JobResponse{
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

// UpdateJob godoc
// @Summary      Update a job
// @Description  Update a job by ID (partial update supported)
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Param        job  body      UpdateJobRequest  true  "Job update info"
// @Success      200  {object}  JobResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /jobs/{id} [put]
func (h *JobHandler) UpdateJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job id"})
		return
	}

	var req UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job"})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Company != nil {
		updates["company"] = *req.Company
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.State != nil {
		updates["state"] = *req.State
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := h.repo.Update(c.Request.Context(), uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
		return
	}

	updatedJob, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated job"})
		return
	}

	c.JSON(http.StatusOK, JobResponse{
		ID:          updatedJob.ID,
		Title:       updatedJob.Title,
		Description: updatedJob.Description,
		Company:     updatedJob.Company,
		City:        updatedJob.City,
		State:       updatedJob.State,
		CreatedAt:   updatedJob.CreatedAt,
		Status:      updatedJob.Status,
	})
}
