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

type ResourceResponse struct {
    Name string `json:"name"`
    RequestCount int `json:"request_count"`
	TimeFrame int `json:"time_frame"`
}

func (s *Server) listResources(w http.ResponseWriter, r *http.Request) {
    resources := make([]ResourceResponse, 0, len(s.resources))
    for _, resource := range s.resources {
        resources = append(resources, ResourceResponse{
            Name: resource.Name,
            RequestCount: resource.RequestCount,
            TimeFrame: resource.TimeFrame,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resources)
}

func (s *Server) updateResource(w http.ResponseWriter, r *http.Request) {
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

    // Update the resource
    if _, exists := s.resources[data.Name]; !exists {
        http.Error(w, "Resource not found", http.StatusNotFound)
        return
    }

    s.resources[data.Name] = Resource{
        Name:           data.Name,
        RequestCount:   data.RequestCount,
        TimeFrame:      data.TimeFrame,
        ScheduledCalls: s.resources[data.Name].ScheduledCalls,
    }

    w.WriteHeader(http.StatusOK)
    message := fmt.Sprintf("Resource %s updated with limit of %d requests per %d seconds\n", data.Name, data.RequestCount, data.TimeFrame)
    w.Write([]byte(message))
}

func (s *Server) deleteResource(w http.ResponseWriter, r *http.Request) {
    var data struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Name == "" {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Get the resource-specific mutex
    resourceMutex, _ := s.resourceMutexes.LoadOrStore(data.Name, &sync.Mutex{})
    mu := resourceMutex.(*sync.Mutex)
    mu.Lock()
    defer mu.Unlock()

    // Delete the resource
    if _, exists := s.resources[data.Name]; !exists {
        http.Error(w, "Resource not found", http.StatusNotFound)
        return
    }

    delete(s.resources, data.Name)

    w.WriteHeader(http.StatusOK)
    message := fmt.Sprintf("Resource %s deleted\n", data.Name)
    w.Write([]byte(message))
}