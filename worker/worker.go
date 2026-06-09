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

		process(j)
	}
}

func process(j job.Job) {
	fmt.Printf("Processing job %s of type %s\n", j.ID, j.Type)

	switch j.Type {
	case "email":
		handleEmail(j)
	case "resize_image":
		handleResizeImage(j)
	default:
		log.Printf("Unknown job type: %s", j.Type)
	}

	fmt.Printf("Job %s completed\n", j.ID)
}

func handleEmail(j job.Job) {
	// real email logic would go here
	fmt.Printf("Sending email with payload: %s\n", j.Payload)
}

func handleResizeImage(j job.Job) {
	// real image resizing logic would go here
	fmt.Printf("Resizing image with payload: %s\n", j.Payload)
}
