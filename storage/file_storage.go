package storage

import (
	"encoding/json"
	"meter_flow/model"
	"os"
)

type FileStorage struct {
	filepath string
}

func NewFileStorage(filepath string) *FileStorage {
	return &FileStorage{filepath: filepath}
}

func (fs *FileStorage) Save(resources map[string]model.Resource) error {
	// Create a DTO map without ScheduledCalls
	persistentData := make(map[string]ResourceDTO)

	for key, resource := range resources {
		persistentData[key] = ResourceDTO{
			Name:         resource.Name,
			RequestCount: resource.RequestCount,
			TimeFrame:    resource.TimeFrame,
		}
	}

	data, err := json.Marshal(persistentData)
	if err != nil {
		return err
	}
	return os.WriteFile(fs.filepath, data, 0644)
}

func (fs *FileStorage) Load() (map[string]model.Resource, error) {
	data, err := os.ReadFile(fs.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]model.Resource), nil
		}
		return nil, err
	}

	var persistentData map[string]ResourceDTO
	if err := json.Unmarshal(data, &persistentData); err != nil {
		return nil, err
	}

	// Convert back to full Resource objects
	resources := make(map[string]model.Resource)
	for key, dto := range persistentData {
		resources[key] = model.Resource{
			Name:           dto.Name,
			RequestCount:   dto.RequestCount,
			TimeFrame:      dto.TimeFrame,
			ScheduledCalls: []int64{}, // Empty slice for scheduled calls
		}
	}

	return resources, nil
}
