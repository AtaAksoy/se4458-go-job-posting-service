package jobs

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type JobRepository interface {
	Create(ctx context.Context, job *Job) error
	List(ctx context.Context, offset, limit int) ([]Job, int64, error)
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string, offset, limit int) ([]Job, int64, error)
	GetByID(ctx context.Context, id uint) (*Job, error)
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
}

type GormJobRepository struct {
	db    *gorm.DB
	cache *JobCache
}

func NewGormJobRepository(db *gorm.DB, cache *JobCache) JobRepository {
	return &GormJobRepository{db: db, cache: cache}
}

func (r *GormJobRepository) Create(ctx context.Context, job *Job) error {
	err := r.db.Create(job).Error
	if err != nil {
		return err
	}

	r.cache.SetJob(ctx, job)

	r.cache.InvalidateJobsList(ctx)

	return nil
}

func (r *GormJobRepository) List(ctx context.Context, offset, limit int) ([]Job, int64, error) {
	page := (offset / limit) + 1

	// Try to get from cache first
	jobs, total, err := r.cache.GetJobsList(ctx, page, limit)
	fmt.Println(jobs)
	if err == nil {
		return jobs, total, nil
	}

	// If not in cache, get from database
	var dbJobs []Job
	var dbTotal int64
	r.db.Model(&Job{}).Count(&dbTotal)
	err = r.db.Order("created_at desc").Offset(offset).Limit(limit).Find(&dbJobs).Error
	if err != nil {
		return nil, 0, err
	}

	r.cache.SetJobsList(ctx, page, limit, dbJobs, dbTotal)

	return dbJobs, dbTotal, nil
}

func (r *GormJobRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.Delete(&Job{}, id).Error
	if err != nil {
		return err
	}

	// Invalidate caches
	r.cache.InvalidateJob(ctx, id)
	r.cache.InvalidateJobsList(ctx)
	r.cache.InvalidateJobsSearch(ctx)

	return nil
}

func (r *GormJobRepository) Search(ctx context.Context, query string, offset, limit int) ([]Job, int64, error) {
	page := (offset / limit) + 1

	// Try to get from cache first
	jobs, total, err := r.cache.GetJobsSearch(ctx, query, page, limit)
	if err == nil {
		return jobs, total, nil
	}

	// If not in cache, get from database
	var dbJobs []Job
	var dbTotal int64
	q := "%" + query + "%"
	dbq := r.db.Model(&Job{}).Where(
		"title LIKE ? OR description LIKE ? OR company LIKE ? OR city LIKE ? OR state LIKE ?",
		q, q, q, q, q,
	)
	dbq.Count(&dbTotal)
	err = dbq.Order("created_at desc").Offset(offset).Limit(limit).Find(&dbJobs).Error
	if err != nil {
		return nil, 0, err
	}

	// Cache the result
	r.cache.SetJobsSearch(ctx, query, page, limit, dbJobs, dbTotal)

	return dbJobs, dbTotal, nil
}

func (r *GormJobRepository) GetByID(ctx context.Context, id uint) (*Job, error) {
	// Try to get from cache first
	job, err := r.cache.GetJob(ctx, id)
	if err == nil {
		return job, nil
	}

	// If not in cache, get from database
	var dbJob Job
	err = r.db.First(&dbJob, id).Error
	if err != nil {
		return nil, err
	}

	// Cache the job
	r.cache.SetJob(ctx, &dbJob)

	return &dbJob, nil
}

func (r *GormJobRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	err := r.db.Model(&Job{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}

	// Invalidate caches
	r.cache.InvalidateJob(ctx, id)
	r.cache.InvalidateJobsList(ctx)
	r.cache.InvalidateJobsSearch(ctx)

	return nil
}
