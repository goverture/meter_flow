package handlers

import (
	"encoding/json"
	"fmt"
	"meter_flow/model"
	"meter_flow/server"
	"net/http"
	"sync"
)

func RegisterResource(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		resourceMutex, _ := srv.ResourceMutexes.LoadOrStore(data.Name, &sync.Mutex{})
		mu := resourceMutex.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()

		// Register the new resource
		if _, exists := srv.Resources[data.Name]; exists {
			http.Error(w, "Resource already exists", http.StatusConflict)
			return
		}

		srv.Resources[data.Name] = model.Resource{
			Name:         data.Name,
			RequestCount: data.RequestCount,
			TimeFrame:    data.TimeFrame,
		}

		w.WriteHeader(http.StatusCreated)
		message := fmt.Sprintf("Resource %s with limit of %d requests per %d seconds registered\n", data.Name, data.RequestCount, data.TimeFrame)
		w.Write([]byte(message))
	}
}

type ResourceResponse struct {
	Name         string `json:"name"`
	RequestCount int    `json:"request_count"`
	TimeFrame    int    `json:"time_frame"`
}

func ListResources(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resources := make([]ResourceResponse, 0, len(srv.Resources))
		for _, resource := range srv.Resources {
			resources = append(resources, ResourceResponse{
				Name:         resource.Name,
				RequestCount: resource.RequestCount,
				TimeFrame:    resource.TimeFrame,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resources)
	}
}

func UpdateResource(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		resourceMutex, _ := srv.ResourceMutexes.LoadOrStore(data.Name, &sync.Mutex{})
		mu := resourceMutex.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()

		// Update the resource
		if _, exists := srv.Resources[data.Name]; !exists {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}

		srv.Resources[data.Name] = model.Resource{
			Name:           data.Name,
			RequestCount:   data.RequestCount,
			TimeFrame:      data.TimeFrame,
			ScheduledCalls: srv.Resources[data.Name].ScheduledCalls,
		}

		w.WriteHeader(http.StatusOK)
		message := fmt.Sprintf("Resource %s updated with limit of %d requests per %d seconds\n", data.Name, data.RequestCount, data.TimeFrame)
		w.Write([]byte(message))
	}
}

func DeleteResource(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Name == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Get the resource-specific mutex
		resourceMutex, _ := srv.ResourceMutexes.LoadOrStore(data.Name, &sync.Mutex{})
		mu := resourceMutex.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()

		// Delete the resource
		if _, exists := srv.Resources[data.Name]; !exists {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}

		delete(srv.Resources, data.Name)

		w.WriteHeader(http.StatusOK)
		message := fmt.Sprintf("Resource %s deleted\n", data.Name)
		w.Write([]byte(message))
	}
}
