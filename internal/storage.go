package storage

import (
	"encoding/json"
	"os"
)

func Load(filepath string, model any) error {
	data, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}

	return nil
}

func Save(filepath string, data any) error {
	marshalData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, marshalData, 0666)
}
