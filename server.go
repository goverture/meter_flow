package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Resource struct {
	Name        string
	RequestCount int // Maximum requests allowed
	TimeFrame   int  // Time frame in seconds
}

type Server struct {
	mu        sync.Mutex
	resources map[string]Resource
}

func NewServer() *Server {
	return &Server{
		resources: make(map[string]Resource),
	}
}

func (s *Server) registerResource(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name         string `json:"name"`
		RequestCount int    `json:"request_count"`
		TimeFrame    int    `json:"time_frame"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.RequestCount <= 0 || data.TimeFrame <= 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Register the new resource
	if _, exists := s.resources[data.Name]; exists {
		http.Error(w, "Resource already exists", http.StatusConflict)
		return
	}

	s.resources[data.Name] = Resource{
		Name:         data.Name,
		RequestCount: data.RequestCount,
		TimeFrame:    data.TimeFrame,
	}

	w.WriteHeader(http.StatusCreated)
	message := fmt.Sprintf("Resource %s with limit of %d requests per %d seconds registered\n", data.Name, data.RequestCount, data.TimeFrame)
	w.Write([]byte(message))
}

func (s *Server) scheduleCalls(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ResourceName string `json:"resource_name"`
		NumCalls     int    `json:"num_calls"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.NumCalls <= 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	resource, exists := s.resources[data.ResourceName]
	s.mu.Unlock()
	if !exists {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Calculate the number of full batches and the remaining calls
	totalCalls := data.NumCalls
	batchSize := resource.RequestCount
	timeFrame := resource.TimeFrame
	var delays []int
	batchCount := totalCalls / batchSize
	remainingCalls := totalCalls % batchSize

	// Schedule full batches with a delay after each batch
	for i := 0; i < batchCount; i++ {
		for j := 0; j < batchSize; j++ {
			delays = append(delays, i*timeFrame) // Each call in the batch gets the same delay
		}
	}

	// Schedule any remaining calls with one additional delay
	if remainingCalls > 0 {
		finalDelay := batchCount * timeFrame
		for k := 0; k < remainingCalls; k++ {
			delays = append(delays, finalDelay)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delays)
}

func main() {
	server := NewServer()

	http.HandleFunc("/resources", server.registerResource)
	http.HandleFunc("/schedule", server.scheduleCalls)

	port := ":8080"
	log.Println("MeterFlow server is running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
