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
				Value:   "/home/data/alac",
				Usage:   "root directory for ALAC files",
			},
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"n"},
				Usage:   "actually save the files",
			},
		},
		Action: func(cCtx *cli.Context) error {
			sa.Dryrun(cCtx.Bool("dryrun"))

			dir := cCtx.String("dir")
			if err := sa.WalkTree(dir); err != nil {
				log.Panic(err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
