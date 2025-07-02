package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AtaAksoy/se4458-go-job-posting-service/internal/v1/db"
)

type JobCache struct {
	redis *db.RedisClient
}

func NewJobCache(redis *db.RedisClient) *JobCache {
	return &JobCache{redis: redis}
}

func (c *JobCache) getJobKey(id uint) string {
	return fmt.Sprintf("job:%d", id)
}

func (c *JobCache) getJobsListKey(page, limit int) string {
	return fmt.Sprintf("jobs:list:%d:%d", page, limit)
}

func (c *JobCache) getJobsSearchKey(query string, page, limit int) string {
	return fmt.Sprintf("jobs:search:%s:%d:%d", query, page, limit)
}

func (c *JobCache) SetJob(ctx context.Context, job *Job) error {
	key := c.getJobKey(job.ID)
	return c.redis.Set(ctx, key, job, 30*time.Minute)
}

func (c *JobCache) GetJob(ctx context.Context, id uint) (*Job, error) {
	key := c.getJobKey(id)
	var job Job
	err := c.redis.Get(ctx, key, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (c *JobCache) SetJobsList(ctx context.Context, page, limit int, jobs []Job, total int64) error {
	key := c.getJobsListKey(page, limit)
	data := map[string]interface{}{
		"jobs":  jobs,
		"total": total,
	}
	return c.redis.Set(ctx, key, data, 15*time.Minute)
}

func (c *JobCache) GetJobsList(ctx context.Context, page, limit int) ([]Job, int64, error) {
	key := c.getJobsListKey(page, limit)
	var data map[string]interface{}
	err := c.redis.Get(ctx, key, &data)
	if err != nil {
		return nil, 0, err
	}

	jobsData := data["jobs"].([]interface{})
	jobs := make([]Job, len(jobsData))
	for i, jobData := range jobsData {
		jobBytes, _ := json.Marshal(jobData)
		json.Unmarshal(jobBytes, &jobs[i])
	}

	total := int64(data["total"].(float64))
	return jobs, total, nil
}

func (c *JobCache) SetJobsSearch(ctx context.Context, query string, page, limit int, jobs []Job, total int64) error {
	key := c.getJobsSearchKey(query, page, limit)
	data := map[string]interface{}{
		"jobs":  jobs,
		"total": total,
	}
	return c.redis.Set(ctx, key, data, 10*time.Minute)
}

func (c *JobCache) GetJobsSearch(ctx context.Context, query string, page, limit int) ([]Job, int64, error) {
	key := c.getJobsSearchKey(query, page, limit)
	var data map[string]interface{}
	err := c.redis.Get(ctx, key, &data)
	if err != nil {
		return nil, 0, err
	}

	jobsData := data["jobs"].([]interface{})
	jobs := make([]Job, len(jobsData))
	for i, jobData := range jobsData {
		jobBytes, _ := json.Marshal(jobData)
		json.Unmarshal(jobBytes, &jobs[i])
	}

	total := int64(data["total"].(float64))
	return jobs, total, nil
}

func (c *JobCache) InvalidateJob(ctx context.Context, id uint) error {
	key := c.getJobKey(id)
	return c.redis.Del(ctx, key)
}

func (c *JobCache) InvalidateJobsList(ctx context.Context) error {
	return c.redis.DelPattern(ctx, "jobs:list:*")
}

func (c *JobCache) InvalidateJobsSearch(ctx context.Context) error {
	return c.redis.DelPattern(ctx, "jobs:search:*")
}

func (c *JobCache) InvalidateAll(ctx context.Context) error {
	return c.redis.DelPattern(ctx, "job:*")
}
