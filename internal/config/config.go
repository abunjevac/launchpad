package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultShortcut = "<Super>Escape"

// Config represents the launchpad configuration.
type Config struct {
	Shortcut string  `yaml:"shortcut,omitempty"`
	Entries  []Entry `yaml:"entries"`
}

// Entry represents a single launchpad entry.
type Entry struct {
	Name      string `yaml:"name,omitempty"`
	Icon      string `yaml:"icon,omitempty"`
	Command   string `yaml:"command,omitempty"`
	Separator bool   `yaml:"separator,omitempty"`
}

// DefaultConfigPath returns the default config file path (~/.launchpad.yaml).
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.Join(home, ".launchpad.yaml"), nil
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	cfg.applyDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config %s: %w", path, err)
	}

	return &cfg, nil
}

func (cfg *Config) applyDefaults() {
	if strings.TrimSpace(cfg.Shortcut) == "" {
		cfg.Shortcut = DefaultShortcut
	}
}

// Validate rejects malformed entries before the UI tries to render or launch
// them.
func (cfg *Config) Validate() error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	for i, entry := range cfg.Entries {
		if entry.Separator {
			continue
		}

		if strings.TrimSpace(entry.Name) == "" {
			return fmt.Errorf("entry %d: name is required", i+1)
		}

		if strings.TrimSpace(entry.Command) == "" {
			return fmt.Errorf("entry %d: command is required", i+1)
		}
	}

	return nil
}
