package sa

import (
	"log"
	"sync"

	"go.uber.org/ratelimit"
)

type Config struct {
	DryRun      bool
	Debug       bool
	SkipMB      bool
	SkipMove    bool
	Overwrite   bool
	FinalDir    string
	AcoustIDKey string
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

	mb5query mb5_query
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
		rl: ratelimit.New(1),
	}

	if err := c.initMB5(); err != nil {
		log.Printf("warning: mb5 init failed: %v", err)
	} else {
		c.mb5query = mb5_query_new("SmartAlac", "musicbrainz.org", 0)
	}

	return c
}

func (c *Curator) ShowStats() {
	c.Stats.mu.Lock()
	defer c.Stats.mu.Unlock()
	log.Printf("%d files / %d changes / %d bad queries", c.Stats.Files, c.Stats.Changes, len(c.Stats.BadQueries))
}
