package server

import (
	"sync"

	"meter_flow/model"
	"meter_flow/storage"
)

type Server struct {
	ResourceMutexes sync.Map // Map of resource name to resource-specific mutex
	Resources       map[string]model.Resource
	storage         storage.Storage
}

func NewServer(storage storage.Storage) *Server {
	// Load initial state
	resources, err := storage.Load()
	if err != nil {
		println("Error loading resources:", err)
		resources = make(map[string]model.Resource)
	}

	return &Server{
		Resources: resources,
		storage:   storage,
	}
}

func (s *Server) Persist() error {
	return s.storage.Save(s.Resources)
}
