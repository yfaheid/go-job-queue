package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/yfaheid/go-job-queue/producer"
)

type Server struct {
	rdb *redis.Client
}

func NewServer(rdb *redis.Client) *Server {
	return &Server{rdb: rdb}
}

func (s *Server) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", s.handleEnqueue)
	mux.HandleFunc("GET /jobs/{id}", s.handleGetJob)
	mux.HandleFunc("GET /jobs", s.handleListJobs)
	return mux
}

type enqueueRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (s *Server) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	var req enqueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := producer.Enqueue(s.rdb, req.Type, req.Payload); err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "enqueued"})
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := context.Background()

	data, err := s.rdb.HGetAll(ctx, "job:"+id).Result()
	if err != nil || len(data) == 0 {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// get the 20 most recent job IDs
	ids, err := s.rdb.LRange(ctx, "job_ids", 0, 19).Result()
	if err != nil {
		http.Error(w, "failed to fetch job ids", http.StatusInternalServerError)
		return
	}

	jobs := []map[string]string{}
	for _, id := range ids {
		data, err := s.rdb.HGetAll(ctx, "job:"+id).Result()
		if err != nil || len(data) == 0 {
			continue
		}
		jobs = append(jobs, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
