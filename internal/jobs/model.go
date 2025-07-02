package jobs

type Job struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	City        string `json:"city"`
	State       string `json:"state"`
	Status      bool   `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

func (Job) TableName() string {
	return "jobs"
}
