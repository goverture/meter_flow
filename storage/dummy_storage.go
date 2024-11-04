package storage

import "meter_flow/model"

// Dummy storage implementation for tests
type DummyStorage struct {
	Resources map[string]model.Resource
}

func NewDummyStorage() *DummyStorage {
	return &DummyStorage{Resources: make(map[string]model.Resource)}
}

func (ds *DummyStorage) Save(resources map[string]model.Resource) error {
	ds.Resources = resources
	return nil
}

func (ds *DummyStorage) Load() (map[string]model.Resource, error) {
	return ds.Resources, nil
}
