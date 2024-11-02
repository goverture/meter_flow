package persistence

import (
	"encoding/json"
	"meter_flow/model"
	"os"
)

type FilePersistence struct {
    filename string
}

func NewFilePersistence(filename string) *FilePersistence {
    return &FilePersistence{filename: filename}
}

func (f *FilePersistence) Save(data *map[string]model.Resource) error {
    file, err := os.Create(f.filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(data)
    if err != nil {
        return err
    }

    return nil
}

func (f *FilePersistence) Load(data *map[string]model.Resource) error {
    file, err := os.Open(f.filename)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // File doesn't exist, no resources to load
        }
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(data)
    if err != nil && err.Error() != "EOF" {
        return err
    }

    return nil
}
