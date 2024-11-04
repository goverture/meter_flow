package main

import (
	"log"
	"meter_flow/handlers"
	"meter_flow/server"
	"meter_flow/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func handleShutdown(server *server.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)

		if err := server.Persist(); err != nil {
			log.Printf("Error saving resources: %v", err)
		} else {
			log.Println("Resources saved successfully")
		}

		os.Exit(0)
	}()
}

func main() {
	storage := storage.NewFileStorage("resources.json")
	server := server.NewServer(storage)
	// save the resources to disk upon shutdown
	handleShutdown(server)

	// "resources" endpoints
	http.HandleFunc("POST /resources", handlers.RegisterResource(server))
	http.HandleFunc("GET /resources", handlers.ListResources(server))
	http.HandleFunc("PUT /resources", handlers.UpdateResource(server))
	http.HandleFunc("DELETE /resources", handlers.DeleteResource(server))

	// "schedule" endpoint
	http.HandleFunc("POST /schedule", handlers.ScheduleCalls(server))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("MeterFlow server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
