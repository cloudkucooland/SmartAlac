package sa

import (
	"context"
	"log"
	"sync"

	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
	"github.com/irlndts/go-discogs"
	"go.uber.org/ratelimit"
)

type Config struct {
	DryRun       bool
	Debug        bool
	SkipMB       bool
	SkipMove     bool
	Overwrite    bool
	FinalDir     string
	AcoustIDKey  string
	DiscogsToken string
	FpcalcPath   string
}

type Stats struct {
	Changes    int
	Files      int
	BadQueries map[string]bool
	mu         sync.Mutex
}

type Curator struct {
	Config Config
	Stats  Stats
	rl     ratelimit.Limiter
	dgRL   ratelimit.Limiter
	Cache  *Cache

	mb5query mb5.Query
	dgClient interface{} // github.com/irlndts/go-discogs
	ctx      context.Context
}

func NewCurator(cfg Config) *Curator {
	if cfg.FinalDir == "" {
		cfg.FinalDir = "/home/music/alac"
	}

	c := &Curator{
		Config: cfg,
		Stats: Stats{
			BadQueries: make(map[string]bool),
		},
		rl:   ratelimit.New(1),
		dgRL: ratelimit.New(1),
	}

	if cfg.DiscogsToken != "" {
		d, err := discogs.New(&discogs.Options{
			Token:     cfg.DiscogsToken,
			UserAgent: "SmartAlac/0.6.0 +https://github.com/cloudkucooland/SmartAlac",
		})
		if err != nil {
			log.Printf("warning: discogs client init failed: %v", err)
		} else {
			c.dgClient = d
		}
	}

	if err := c.initMB5(); err != nil {
		log.Printf("warning: mb5 init failed: %v", err)
	} else {
		c.mb5query = mb5.QueryNew("SmartAlac", "musicbrainz.org", 0)
	}

	return c
}

func (c *Curator) ShowStats() {
	c.Stats.mu.Lock()
	defer c.Stats.mu.Unlock()
	log.Printf("%d files / %d changes / %d bad queries", c.Stats.Files, c.Stats.Changes, len(c.Stats.BadQueries))
}
