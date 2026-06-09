package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/yfaheid/go-job-queue/api"
	"github.com/yfaheid/go-job-queue/worker"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	go worker.Start(rdb)

	server := api.NewServer(rdb)
	mux := server.RegisterRoutes()

	fmt.Println("API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
