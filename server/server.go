package server

import (
	"sync"

	"meter_flow/model"
)

type Server struct {
	ResourceMutexes sync.Map // Map of resource name to resource-specific mutex
	Resources       map[string]model.Resource
}

func NewServer() *Server {
	return &Server{
		Resources: make(map[string]model.Resource),
	}
}