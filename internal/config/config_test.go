package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsEntryWithoutName(t *testing.T) {
	path := writeConfig(t, `
entries:
  - command: mate-terminal
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	if !strings.Contains(err.Error(), "entry 1: name is required") {
		t.Fatalf("Load() error = %q, want missing name", err)
	}
}

func TestLoadRejectsEntryWithoutCommand(t *testing.T) {
	path := writeConfig(t, `
entries:
  - name: Terminal
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	if !strings.Contains(err.Error(), "entry 1: command is required") {
		t.Fatalf("Load() error = %q, want missing command", err)
	}
}

func TestLoadAllowsSeparatorsWithoutNameOrCommand(t *testing.T) {
	path := writeConfig(t, `
entries:
  - separator: true
  - name: Terminal
    command: mate-terminal
`)

	if _, err := Load(path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "launchpad.yaml")

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return path
}
