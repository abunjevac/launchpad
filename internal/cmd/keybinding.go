package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	commandSlot = "command-999"
	runSlot     = "run-command-999"
)

// RegisterKeybinding registers the launchpad shortcut via Mate's Marco gsettings.
func RegisterKeybinding(binaryPath, shortcut string) error {
	launchCommand := shellQuote(binaryPath) + " run"

	regCmd := exec.Command(
		"gsettings",
		"set",
		"org.mate.Marco.keybinding-commands",
		commandSlot,
		launchCommand,
	)
	if out, err := regCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("setting command %s: %w\nOutput: %s", commandSlot, err, string(out))
	}

	regKey := exec.Command(
		"gsettings",
		"set",
		"org.mate.Marco.global-keybindings",
		runSlot,
		shortcut,
	)
	if out, err := regKey.CombinedOutput(); err != nil {
		return fmt.Errorf("setting keybinding %s: %w\nOutput: %s", runSlot, err, string(out))
	}

	fmt.Printf("Registered shortcut %s for %s\n", shortcut, launchCommand)

	return nil
}

// UnregisterKeybinding removes the launchpad shortcut from Mate's Marco.
func UnregisterKeybinding() error {
	unregCmd := exec.Command(
		"gsettings",
		"reset",
		"org.mate.Marco.keybinding-commands",
		commandSlot,
	)
	if out, err := unregCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("resetting command %s: %w\nOutput: %s", commandSlot, err, string(out))
	}

	unregKey := exec.Command(
		"gsettings",
		"reset",
		"org.mate.Marco.global-keybindings",
		runSlot,
	)
	if out, err := unregKey.CombinedOutput(); err != nil {
		return fmt.Errorf("resetting keybinding %s: %w\nOutput: %s", runSlot, err, string(out))
	}

	fmt.Println("Unregistered launchpad shortcut")

	return nil
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}

	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
