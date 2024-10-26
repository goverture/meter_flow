package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	mu        sync.Mutex
	resources map[string]bool
}

func NewServer() *Server {
	return &Server{
		resources: make(map[string]bool),
	}
}

func (s *Server) registerResource(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
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

	s.resources[data.Name] = true

	w.WriteHeader(http.StatusCreated)
	message := fmt.Sprintf("Resource %s registered\n", data.Name)
    w.Write([]byte(message))
}

func (s *Server) scheduleCalls(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ResourceName string `json:"resource_name"`
		NumCalls     int    `json:"num_calls"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if _, exists := s.resources[data.ResourceName]; !exists {
		s.mu.Unlock()
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}
	s.mu.Unlock()

	// Assume the resource can handle 100 requests per second
	maxRate := 100
	totalCalls := data.NumCalls

	// Calculate delays in whole seconds
	var delays []int
	for i := 0; i < totalCalls; i++ {
		delay := i / maxRate // Divide to get the second in which the call can be made
		delays = append(delays, delay)
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
