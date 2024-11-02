package main

import (
	"log"
	"net/http"
	"os"
	"sync"
)

type Resource struct {
	Name           string
	RequestCount   int     // Maximum requests allowed
	TimeFrame      int     // Time frame in seconds
	ScheduledCalls []int64 // Track scheduled timestamps for this resource
}

type Server struct {
	resourceMutexes sync.Map // Map of resource name to resource-specific mutex
	resources       map[string]Resource
}

func NewServer() *Server {
	return &Server{
		resources: make(map[string]Resource),
	}
}

func main() {
	server := NewServer()

	// "resources" endpoints
	http.HandleFunc("POST /resources", server.registerResource)
	http.HandleFunc("GET /resources", server.listResources)
	http.HandleFunc("PUT /resources", server.updateResource)
	http.HandleFunc("DELETE /resources", server.deleteResource)

	// "schedule" endpoint
	http.HandleFunc("POST /schedule", server.scheduleCalls)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Println("MeterFlow server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
