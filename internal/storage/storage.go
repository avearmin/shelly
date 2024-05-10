package storage

import (
	"encoding/json"
	"os"
)

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Action      string `json:"Action"`
}

func LoadCommands(filepath string) (map[string]Command, error) {
	commands := make(map[string]Command)

	data, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		return commands, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &commands); err != nil {
		return nil, err
	}

	return commands, nil
}

func SaveCommands(filepath string, commands map[string]Command) error {
	data, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0666)
}
