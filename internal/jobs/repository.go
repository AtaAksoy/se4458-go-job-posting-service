package jobs

import (
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(job *Job) error
	List(offset, limit int) ([]Job, int64, error)
	Delete(id uint) error
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
