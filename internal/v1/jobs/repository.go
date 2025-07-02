package jobs

import (
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(job *Job) error
	List(offset, limit int) ([]Job, int64, error)
	Delete(id uint) error
	Search(query string, offset, limit int) ([]Job, int64, error)
}

type GormJobRepository struct {
	db *gorm.DB
}

func NewGormJobRepository(db *gorm.DB) JobRepository {
	return &GormJobRepository{db: db}
}

func (r *GormJobRepository) Create(job *Job) error {
	return r.db.Create(job).Error
}

func (r *GormJobRepository) List(offset, limit int) ([]Job, int64, error) {
	var jobs []Job
	var total int64
	r.db.Model(&Job{}).Count(&total)
	err := r.db.Order("created_at desc").Offset(offset).Limit(limit).Find(&jobs).Error
	return jobs, total, err
}

func (r *GormJobRepository) Delete(id uint) error {
	return r.db.Delete(&Job{}, id).Error
}

func (r *GormJobRepository) Search(query string, offset, limit int) ([]Job, int64, error) {
	var jobs []Job
	var total int64
	q := "%" + query + "%"
	dbq := r.db.Model(&Job{}).Where(
		"title LIKE ? OR description LIKE ? OR company LIKE ? OR city LIKE ? OR state LIKE ?",
		q, q, q, q, q,
	)
	dbq.Count(&total)
	err := dbq.Order("created_at desc").Offset(offset).Limit(limit).Find(&jobs).Error
	return jobs, total, err
}
