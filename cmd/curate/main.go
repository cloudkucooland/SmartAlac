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
		Version: "v0.0.0",
		Authors: []*cli.Author{
			{
				Name:  "Scot C. Bontrager",
				Email: "cloudkucooland@gmail.com",
			},
		},
		Copyright: "© 2022 Scot C. Bontrager",
		HelpName:  "curate",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   "/home/music/alac",
				Usage:   "root directory for ALAC files",
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
		},
		Action: func(cCtx *cli.Context) error {
			sa.Dryrun(cCtx.Bool("dryrun"))
			sa.Debug(cCtx.Bool("debug"))

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
