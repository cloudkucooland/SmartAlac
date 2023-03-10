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
				// Value:   "/home/data/alac",
				// Value:   "/home/data/alac/xx, The",
				Value:   "/home/data/alac/Material Issue/What Girls Want [1992,CD,US,Mercury Records,CDP 685]/",
				Usage:   "root directory for ALAC files",
			},
		},
        Action: func(cCtx *cli.Context) error {
            dir:= cCtx.String("dir")
            err := sa.WalkTree(dir)
            if err != nil {
                log.Panic(err)
            }
            err = sa.Rename()
            if err != nil {
                log.Panic(err)
            }
            return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
