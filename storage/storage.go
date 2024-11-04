package storage

import "meter_flow/model"

// "Data Transfer Object" for resources, we don't want to store the "ScheduledCalls"
type ResourceDTO struct {
	Name         string
	RequestCount int
	TimeFrame    int
}

// Store and load the server data (resources)
type Storage interface {
	Save(resources map[string]model.Resource) error
	Load() (map[string]model.Resource, error)
}
