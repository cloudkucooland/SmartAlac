package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudkucooland/SmartAlac/pkg/sa"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:      "curate",
		Usage:     "Smart ALAC curator using MusicBrainz",
		Version:   "v0.6.0",
		Authors:   []any{"Scot C. Bontrager <cloudkucooland@gmail.com>"},
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
			&cli.StringFlag{
				Name:    "acoustid-key",
				Usage:   "AcoustID API key for fingerprinting fallback",
				Sources: cli.EnvVars("ACOUSTID_KEY"),
			},
			&cli.StringFlag{
			        Name:  "discogs-token",
			        Usage: "Discogs API token",
			},
			&cli.StringFlag{
			        Name:  "fpcalc",
			        Usage: "path to fpcalc binary",
			        Value: "fpcalc",
			},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {

			shared := sa.LoadSharedConfig()

			cfg := sa.Config{
			        DryRun:       cmd.Bool("dryrun"),
			        Debug:        cmd.Bool("debug"),
			        SkipMB:       cmd.Bool("skipmb"),
			        SkipMove:     cmd.Bool("skipmove"),
			        Overwrite:    cmd.Bool("overwrite"),
			        FinalDir:     cmd.String("finaldir"),
			        AcoustIDKey:  cmd.String("acoustid-key"),
			        DiscogsToken: cmd.String("discogs-token"),
			        FpcalcPath:   cmd.String("fpcalc"),
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
			if cfg.FpcalcPath == "fpcalc" && shared.FpcalcPath != "" {
			        cfg.FpcalcPath = shared.FpcalcPath
			}

			curator := sa.NewCurator(cfg)

			dir := cmd.String("dir")
			if err := curator.WalkTree(ctx, dir); err != nil {
			        return err
			}

			curator.ShowStats()
			return nil
			},

	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
