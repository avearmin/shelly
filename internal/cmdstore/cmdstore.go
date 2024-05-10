package cmdstore

import (
	"os"
	"github.com/avearmin/shelly/internal/storage"
)

type Command struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Action string `json:"action"`
}

func Load(filepath string) (map[string]Command, error) {
	cmds := make(map[string]Command)

	if err := storage.Load(filepath, cmds); err != nil {
		return nil, err
	}

	return cmds, nil
}

func Save(filepath string, cmds map[string]Command) error {
	return storage.Save(filepath, cmds)
}

func GetDefaultPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.config/shelly/commands.json"
	}

	return homeDir + "/.config/shelly/commands.json"
}
