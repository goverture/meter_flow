package handlers

import (
	"encoding/json"
	"meter_flow/scheduler"
	"meter_flow/server"
	"net/http"
	"sync"
	"time"
)

func ScheduleCalls(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			ResourceName string `json:"resource_name"`
			NumCalls     int    `json:"num_calls"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.NumCalls <= 0 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Get the resource-specific mutex
		resourceMutex, _ := srv.ResourceMutexes.LoadOrStore(data.ResourceName, &sync.Mutex{})
		mu := resourceMutex.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()

		resource, exists := srv.Resources[data.ResourceName]
		if !exists {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}

		// Get the current time and schedule new calls
		now := time.Now().Unix()
		delays, updatedCalls := scheduler.Schedule(data.NumCalls, resource.RequestCount, resource.TimeFrame, resource.ScheduledCalls, now)

		// Update the resource with the latest scheduled calls
		resource.ScheduledCalls = updatedCalls
		srv.Resources[data.ResourceName] = resource

		response := map[string]interface{}{
			"delays": delays,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
