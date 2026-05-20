package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/abunjevac/launchpad/internal/config"
)

// InitConfig generates a skeleton config file at the specified path.
func InitConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config already exists at %s", path)
	}

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	cfg := config.Config{
		Entries: []config.Entry{
			{
				Name:    "Terminal",
				Icon:    "utilities-terminal",
				Command: "mate-terminal",
			},
			{
				Separator: true,
			},
			{
				Name:    "Files",
				Icon:    "system-file-manager",
				Command: "caja .",
			},
			{
				Name:    "Text Editor",
				Icon:    "accessories-text-editor",
				Command: "pluma %U",
			},
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config to %s: %w", path, err)
	}

	fmt.Printf("Created skeleton config at %s\n", path)

	return nil
}
