package configstore

import (
	"os"
	"github.com/avearmin/shelly/internal/storage"
)

type Config struct {
	CmdsPath string `json:"cmdspath"`
}

func Load() (Config, error) {
	config := Config{}

	if err := storage.Load(GetPath(), config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func Save(config Config) error {
	return storage.Save(GetPath(), config)
}

func GetPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.config/shelly/config.json"
	}

	return homeDir + "/.config/shelly/config.json"
}

func Create() error {
	_, err := os.Create(GetPath())
	return err
}
