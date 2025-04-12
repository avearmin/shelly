package cmdstore

import (
	"fmt"
	"os"
	"time"

	"github.com/avearmin/shelly/internal/storage"
)

type Command struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LastUsed    time.Time `json:"last_used"`
	Action      string    `json:"action"`
}

func (c Command) LastUsedInHumanTerms() string {
	if c.LastUsed.IsZero() {
		return "Never"
	}

	now := time.Now()
	diff := now.Sub(c.LastUsed)

	days := int(diff.Hours() / 24)

	switch {
	case days < 1:
		return "Today"
	case days == 1:
		return "Yesterday"
	case days < 7:
		return fmt.Sprintf("%d days ago", days)
	case days < 30:
		weeks := days / 7
		return fmt.Sprintf("%d week%s ago", weeks, plural(weeks))
	case days < 365:
		months := days / 30
		return fmt.Sprintf("%d month%s ago", months, plural(months))
	default:
		years := days / 365
		return fmt.Sprintf("%d year%s ago", years, plural(years))
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func Load(filepath string) (map[string]Command, error) {
	cmds := make(map[string]Command)

	if err := storage.Load(filepath, &cmds); err != nil {
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
