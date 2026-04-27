package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/cloudkucooland/SmartAlac/pkg/sa"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "missing",
		Usage:   "Identify folders missing tracks based on MusicBrainz cache",
		Version: "v0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "directory to process",
				Value:   "/home/music/alac",
			},
			&cli.StringFlag{
				Name:    "cache",
				Aliases: []string{"c"},
				Usage:   "path to sqlite cache",
				Value:   sa.DefaultCachePath(),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cache, err := sa.OpenCache(cmd.String("cache"))
			if err != nil {
				return err
			}
			defer cache.Close()

			dir := cmd.String("dir")
			
			// We want to find every directory that contains .m4a files
			// and check if the track count matches the MBID in those files.
			
			albumDirs := make(map[string]bool)
			filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
				if err != nil { return err }
				if !d.IsDir() && strings.HasSuffix(p, ".m4a") {
					albumDirs[filepath.Dir(p)] = true
				}
				return nil
			})

			for albumDir := range albumDirs {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				files, err := os.ReadDir(albumDir)
				if err != nil {
					continue
				}

				var m4aFiles []string
				var mbid string
				for _, f := range files {
					if strings.HasSuffix(f.Name(), ".m4a") {
						m4aFiles = append(m4aFiles, f.Name())
						if mbid == "" {
							// Read MBID from the first file we find
							mp4, err := mp4tag.Open(filepath.Join(albumDir, f.Name()))
							if err == nil {
								tags, err := mp4.Read()
								if err == nil {
									mbid = tags.Custom["MusicBrainz Album Id"]
								}
								mp4.Close()
							}
						}
					}
				}

				if mbid == "" {
					continue
				}

				// Check cache for expected track count
				rm, err := cache.GetRelease(mbid)
				if err != nil {
					// Not in cache, we can't determine if it's missing tracks
					continue
				}

				if len(m4aFiles) < rm.TrackCount {
					fmt.Printf("%s: Missing tracks (have %d, expect %d) [MBID: %s]\n", 
						albumDir, len(m4aFiles), rm.TrackCount, mbid)
				}
			}

			return nil
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
