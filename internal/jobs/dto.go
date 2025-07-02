package jobs

type CreateJobRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Company     string `json:"company" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
}

type JobResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	City        string `json:"city"`
	State       string `json:"state"`
	CreatedAt   int64  `json:"created_at"`
	Status      bool   `json:"status"`
}
