package api

import (
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
