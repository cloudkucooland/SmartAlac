package main

import (
	"log"
	"os"

	"github.com/cloudkucooland/SmartAlac"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "curate",
		Version: "v0.5.0",
		Authors: []*cli.Author{
			{
				Name:  "Scot C. Bontrager",
				Email: "cloudkucooland@gmail.com",
			},
		},
		Copyright: "© 2025 Scot C. Bontrager",
		HelpName:  "curate",

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
				Usage:   "skip polling musicbrainz",
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
		Action: func(cCtx *cli.Context) error {
			sa.Dryrun(cCtx.Bool("dryrun"))
			sa.Debug(cCtx.Bool("debug"))
			sa.Overwrite(cCtx.Bool("overwrite"))
			sa.Finaldir(cCtx.String("finaldir"))

			dir := cCtx.String("dir")
			if err := sa.WalkTree(dir); err != nil {
				log.Panic(err)
			}

			sa.ShowStats()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
