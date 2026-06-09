package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yfaheid/go-job-queue/job"
)

func Enqueue(rdb *redis.Client, jobType string, payload string) error {
	j := job.Job{
		ID:          uuid.NewString(),
		Type:        jobType,
		Payload:     payload,
		Status:      "pending",
		Attempts:    0,
		MaxAttempts: 3,
	}

	data, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	ctx := context.Background()
	return rdb.LPush(ctx, "jobs", data).Err()
}
