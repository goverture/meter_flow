package main

import (
	"log"
	"net/http"
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

	http.HandleFunc("POST /resources", server.registerResource)
	http.HandleFunc("POST /schedule", server.scheduleCalls)

	port := ":8080"
	log.Println("MeterFlow server is running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
