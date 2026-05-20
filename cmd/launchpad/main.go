package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/abunjevac/launchpad/internal/cmd"
	"github.com/abunjevac/launchpad/internal/config"
	"github.com/abunjevac/launchpad/internal/ui"
)

func main() {
	app := &cli.Command{
		Name:  "launchpad",
		Usage: "A popup launcher for Ubuntu Mate",
		Flags: []cli.Flag{
			configFlag(),
		},
		Action: runAction,
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the launchpad UI popup",
				Flags: []cli.Flag{
					configFlag(),
				},
				Action: runAction,
			},
			{
				Name:  "register",
				Usage: "Register the global shortcut key (default: <Super>F1)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Path to config file (to read shortcut override)",
					},
					&cli.StringFlag{
						Name:    "shortcut",
						Aliases: []string{"s"},
						Usage:   "Shortcut key combination (overrides config)",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					shortcut := c.String("shortcut")

					if shortcut == "" {
						cfgPath := resolveConfigPath(c.String("config"))

						if cfg, err := config.Load(cfgPath); err == nil {
							shortcut = cfg.Shortcut
						} else {
							shortcut = config.DefaultShortcut
						}
					}

					binaryPath, err := os.Executable()
					if err != nil {
						return fmt.Errorf("getting binary path: %w", err)
					}

					binaryPath, err = filepath.Abs(binaryPath)
					if err != nil {
						return fmt.Errorf("resolving binary path: %w", err)
					}

					return cmd.RegisterKeybinding(binaryPath, shortcut)
				},
			},
			{
				Name:  "unregister",
				Usage: "Remove the global shortcut key",
				Action: func(ctx context.Context, c *cli.Command) error {
					return cmd.UnregisterKeybinding()
				},
			},
			{
				Name:  "init",
				Usage: "Generate a skeleton config file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Path to config file (default: ~/.launchpad.yaml)",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					cfgPath := resolveConfigPath(c.String("config"))

					return cmd.InitConfig(cfgPath)
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}

func configFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Path to config file (default: ~/.launchpad.yaml)",
	}
}

func runAction(ctx context.Context, c *cli.Command) error {
	_ = ctx

	cfgPath := resolveConfigPath(c.String("config"))

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ui.Run(cfg)

	return nil
}

func resolveConfigPath(path string) string {
	if path != "" {
		return path
	}

	p, err := config.DefaultConfigPath()
	if err != nil {
		return ".launchpad.yaml"
	}

	return p
}
