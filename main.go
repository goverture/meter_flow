package main

import (
	"log"
	"meter_flow/handlers"
	"meter_flow/server"
	"net/http"
	"os"
)


func main() {
	server := server.NewServer()

	// "resources" endpoints
	http.HandleFunc("POST /resources", handlers.RegisterResource(server))
	http.HandleFunc("GET /resources", handlers.ListResources(server))
	http.HandleFunc("PUT /resources", handlers.UpdateResource(server))
	http.HandleFunc("DELETE /resources", handlers.DeleteResource(server))

	// "schedule" endpoint
	http.HandleFunc("POST /schedule", handlers.ScheduleCalls(server))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Println("MeterFlow server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
