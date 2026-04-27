package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudkucooland/SmartAlac/pkg/picker"
	"github.com/cloudkucooland/SmartAlac/pkg/sa"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "picker",
		Usage:   "Interactive album picker and tagger",
		Version: "v0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "directory to process",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "finaldir",
				Aliases: []string{"D"},
				Value:   "/home/music/alac",
				Usage:   "where to move files",
			},
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"n"},
				Usage:   "skip saving the files",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"V"},
				Usage:   "verbose info dumps",
			},
			&cli.StringFlag{
				Name:  "acoustid-key",
				Usage: "AcoustID API key",
			},
			&cli.StringFlag{
				Name:  "discogs-token",
				Usage: "Discogs API token",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			shared := sa.LoadSharedConfig()
			dir, err := filepath.Abs(cmd.String("dir"))
			if err != nil {
				return err
			}

			cfg := sa.Config{
				DryRun:       cmd.Bool("dryrun"),
				Debug:        cmd.Bool("debug"),
				FinalDir:     cmd.String("finaldir"),
				AcoustIDKey:  cmd.String("acoustid-key"),
				DiscogsToken: cmd.String("discogs-token"),
			}

			if cfg.FinalDir == "" {
				cfg.FinalDir = shared.FinalDir
			}
			if cfg.AcoustIDKey == "" {
				cfg.AcoustIDKey = shared.AcoustIDKey
			}
			if cfg.DiscogsToken == "" {
				cfg.DiscogsToken = shared.DiscogsToken
			}

			curator := sa.NewCurator(cfg)
			
			p := tea.NewProgram(picker.NewModel(curator, dir), tea.WithAltScreen())
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
