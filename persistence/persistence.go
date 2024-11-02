package persistence

import "meter_flow/model"


type Persistence interface {
    Save(data *map[string]model.Resource) error
    Load(data *map[string]model.Resource) error
}
