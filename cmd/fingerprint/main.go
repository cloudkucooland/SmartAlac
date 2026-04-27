package main

import (
	"context"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/cloudkucooland/SmartAlac/pkg/sa"
	"github.com/jo-hoe/chromaprint"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "fingerprint",
		Usage:   "Bulk calculate AcoustID fingerprints and store in tags/cache",
		Version: "v0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "directory to process",
				Value:   ".",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "force recalculation even if already tagged",
			},
			&cli.StringFlag{
				Name:  "acoustid-key",
				Usage: "AcoustID API key (optional, if you want to resolve to AcoustIDs)",
			},
			&cli.StringFlag{
				Name:  "fpcalc",
				Usage: "path to fpcalc binary",
				Value: "fpcalc",
			},
			&cli.IntFlag{
				Name:    "jobs",
				Aliases: []string{"j"},
				Usage:   "number of parallel hashing jobs",
				Value:   runtime.NumCPU(),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			shared := sa.LoadSharedConfig()

			cfg := sa.Config{
				AcoustIDKey: cmd.String("acoustid-key"),
				FpcalcPath:  cmd.String("fpcalc"),
			}
			if cfg.AcoustIDKey == "" {
				cfg.AcoustIDKey = shared.AcoustIDKey
			}
			if cfg.FpcalcPath == "fpcalc" && shared.FpcalcPath != "" {
				cfg.FpcalcPath = shared.FpcalcPath
			}

			curator := sa.NewCurator(cfg)

			dir := cmd.String("dir")
			force := cmd.Bool("force")
			numWorkers := int(cmd.Int("jobs"))

			// Worker pool setup
			paths := make(chan string, numWorkers*2)
			var wg sync.WaitGroup
			var errOnce sync.Once
			var walkErr error

			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					builder := chromaprint.NewBuilder().WithPathToChromaprint(cfg.FpcalcPath)
					chromaprinter, err := builder.Build()
					if err != nil {
						errOnce.Do(func() { walkErr = err })
						return
					}

					for p := range paths {
						select {
						case <-ctx.Done():
							return
						default:
						}

						// If not in cache, check if it is already tagged
						mp4, err := mp4tag.Open(p)
						if err == nil {
							tags, err := mp4.Read()
							if err == nil {
								if existingFP := tags.Custom["Acoustid Fingerprint"]; existingFP != "" && !force {
									log.Printf("Found existing fingerprint in tags: %s", p)
									mp4.Close()
									continue
								}
							}
							mp4.Close()
						}

						log.Printf("Hashing: %s", p)
						fingerprints, err := chromaprinter.CreateFingerprints(p)
						if err != nil {
							log.Printf("failed to fingerprint %s: %v", p, err)
							continue
						}

						if len(fingerprints) == 0 {
							continue
						}

						f := fingerprints[0]
						fpStr := sa.EncodeFingerprint(f.Fingerprint)

						var resolvedAID string
						if cfg.AcoustIDKey != "" {
							resolvedAID, _ = curator.AcoustIDLookup(ctx, fpStr, int(f.DurationInSeconds))
						}

						_ = syncTags(ctx, p, fpStr, int(f.DurationInSeconds), resolvedAID)
					}
				}()
			}
			walkErr = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() || !strings.HasSuffix(p, ".m4a") {
					return nil
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case paths <- p:
				}
				return nil
			})

			close(paths)
			wg.Wait()
			return walkErr
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func syncTags(ctx context.Context, path, fp string, dur int, aid string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	mp4, err := mp4tag.Open(path)
	if err != nil {
		return err
	}
	defer mp4.Close()
	mp4.UpperCustom(false)

	tags, err := mp4.Read()
	if err != nil {
		return err
	}

	changed := false
	if tags.Custom == nil {
		tags.Custom = make(map[string]string)
	}

	if tags.Custom["Acoustid Id"] != aid && aid != "" {
		tags.Custom["Acoustid Id"] = aid
		changed = true
	}
	if tags.Custom["Acoustid Fingerprint"] != fp {
		tags.Custom["Acoustid Fingerprint"] = fp
		changed = true
	}

	if changed {
		log.Printf("Updating tags for %s", path)
		// We must specify the keys to save in go-mp4tag local fork
		customKeys := make([]string, 0, len(tags.Custom))
		for k := range tags.Custom {
			customKeys = append(customKeys, k)
		}
		return mp4.Write(tags, customKeys)
	}
	return nil
}
