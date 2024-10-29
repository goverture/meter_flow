package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

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

	// Get the resource-specific mutex
	resourceMutex, _ := s.resourceMutexes.LoadOrStore(data.Name, &sync.Mutex{})
	mu := resourceMutex.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

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
