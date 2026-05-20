package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abunjevac/launchpad/internal/version"
	"github.com/urfave/cli/v3"

	"github.com/abunjevac/launchpad/internal/cmd"
	"github.com/abunjevac/launchpad/internal/config"
	"github.com/abunjevac/launchpad/internal/ui"
)

func main() {
	app := &cli.Command{
		Name:    "launchpad",
		Version: version.Version,
		Usage:   "A popup launcher for Ubuntu Mate",
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
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)

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
