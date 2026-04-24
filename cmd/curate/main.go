package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudkucooland/SmartAlac"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "curate",
		Usage:   "Smart ALAC curator using MusicBrainz",
		Version: "v0.6.0",
		Authors: []any{"Scot C. Bontrager <cloudkucooland@gmail.com>"},
		Copyright: "© 2025 Scot C. Bontrager",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   "/home/music/alac",
				Usage:   "directory to process",
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
				Name:    "skipmb",
				Aliases: []string{"S"},
				Usage:   "skip polling musicbrainz",
			},
			&cli.BoolFlag{
				Name:    "skipmove",
				Aliases: []string{"M"},
				Usage:   "skip moving files",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"V"},
				Usage:   "verbose info dumps",
			},
			&cli.BoolFlag{
				Name:    "overwrite",
				Aliases: []string{"O"},
				Usage:   "overwrite files if duplicates exist",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := sa.Config{
				DryRun:    cmd.Bool("dryrun"),
				Debug:     cmd.Bool("debug"),
				SkipMB:    cmd.Bool("skipmb"),
				SkipMove:  cmd.Bool("skipmove"),
				Overwrite: cmd.Bool("overwrite"),
				FinalDir:  cmd.String("finaldir"),
			}

			curator := sa.NewCurator(cfg)
			dir := cmd.String("dir")
			if err := curator.WalkTree(dir); err != nil {
				return err
			}

			curator.ShowStats()
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
