package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudkucooland/SmartAlac/pkg/bme"
	"github.com/cloudkucooland/SmartAlac/pkg/sa"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "bme",
		Version: "v0.5.0",
		Authors: []any{
			"Scot C. Bontrager <cloudkucooland@gmail.com>",
		},
		Copyright: "© 2026 Scot C. Bontrager",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   "/home/data/bme",
				Usage:   "base directory for work files",
			},
			&cli.StringSliceFlag{
				Name:    "dev",
				Aliases: []string{"D"},
				Usage:   "CD devices to use",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to JSON config file",
			},
			&cli.StringFlag{
				Name:    "encoder",
				Aliases: []string{"e"},
				Usage:   "path to alacenc binary",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"V"},
				Usage:   "verbose info dumps",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			bme.Debug(cmd.Bool("debug"))

			cfgPath := cmd.String("config")
			if cfgPath == "" {
				cfgPath = sa.DefaultConfigPath("bme.json")
			}

			if cmd.Bool("debug") {
				log.Printf("Loading config from: %s", cfgPath)
			}

			cfg, err := bme.LoadConfig(cfgPath)
			if err != nil {
				return err
			}

			if cmd.Bool("debug") {
				log.Printf("Config loaded. Devices: %v", cfg.Devices)
			}

			// CLI Overrides
			if d := cmd.String("dir"); d != "" {
				cfg.BaseDir = d
				// Recalculate subdirs if BaseDir is overridden
				cfg.RipDir = filepath.Join(d, "rip")
				cfg.EncodeDir = filepath.Join(d, "encode")
				cfg.TagDir = filepath.Join(d, "tag")
				cfg.DoneDir = filepath.Join(d, "done")
			}
			if e := cmd.String("encoder"); e != "" {
				cfg.EncoderPath = e
			}

			if devs := cmd.StringSlice("dev"); len(devs) > 0 {
				cfg.Devices = devs
			}

			// Catch OS signals for graceful shutdown
			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
			defer cancel()

			p := tea.NewProgram(bme.NewModel(), tea.WithAltScreen())

			// Start logic in a goroutine
			go func() {
				if err := bme.Start(ctx, cfg, p); err != nil {
					p.Send(bme.StatusMsg{Component: "system", Status: err.Error()})
				}
				p.Quit()
			}()

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
