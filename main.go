package main

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/yfaheid/go-job-queue/producer"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := producer.Enqueue(rdb, "email", `{"to": "user@example.com"}`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Job enqueued!")
}
