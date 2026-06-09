package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yfaheid/go-job-queue/job"
)

func Start(rdb *redis.Client) {
	fmt.Println("Worker started, waiting for jobs...")

	for {
		ctx := context.Background()
		result, err := rdb.BRPop(ctx, 5*time.Second, "jobs").Result()
		if err == redis.Nil {
			fmt.Println("No jobs found, waiting...")
			continue
		}
		if err != nil {
			log.Printf("Error popping job: %v", err)
			continue
		}

		var j job.Job
		if err := json.Unmarshal([]byte(result[1]), &j); err != nil {
			log.Printf("Failed to unmarshal job: %v", err)
			continue
		}

		j.Attempts++
		saveStatus(rdb, j, "running")

		if err := process(j); err != nil {
			log.Printf("Job %s failed (attempt %d/%d): %v", j.ID, j.Attempts, j.MaxAttempts, err)
			if j.Attempts < j.MaxAttempts {
				saveStatus(rdb, j, "pending")
				requeue(rdb, j)
			} else {
				saveStatus(rdb, j, "failed")
				log.Printf("Job %s exhausted all retries", j.ID)
			}
			continue
		}

		saveStatus(rdb, j, "completed")
		fmt.Printf("Job %s completed\n", j.ID)
	}
}

func saveStatus(rdb *redis.Client, j job.Job, status string) {
	ctx := context.Background()
	key := "job:" + j.ID
	rdb.HSet(ctx, key, map[string]interface{}{
		"id":           j.ID,
		"type":         j.Type,
		"status":       status,
		"attempts":     j.Attempts,
		"max_attempts": j.MaxAttempts,
	})
}

func requeue(rdb *redis.Client, j job.Job) {
	data, err := json.Marshal(j)
	if err != nil {
		log.Printf("Failed to marshal job for requeue: %v", err)
		return
	}
	ctx := context.Background()
	rdb.LPush(ctx, "jobs", data)
}

func process(j job.Job) error {
	fmt.Printf("Processing job %s of type %s (attempt %d/%d)\n", j.ID, j.Type, j.Attempts, j.MaxAttempts)

	switch j.Type {
	case "email":
		return handleEmail(j)
	case "resize_image":
		return handleResizeImage(j)
	default:
		return fmt.Errorf("unknown job type: %s", j.Type)
	}
}

func handleEmail(j job.Job) error {
	fmt.Printf("Sending email with payload: %s\n", j.Payload)
	return nil
}

func handleResizeImage(j job.Job) error {
	fmt.Printf("Resizing image with payload: %s\n", j.Payload)
	return nil
}
