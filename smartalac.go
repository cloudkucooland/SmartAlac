package sa

import (
	"log"
)

var dryrun bool
var debug bool

func Dryrun(b bool) {
	dryrun = b
}

func Debug(b bool) {
	debug = b
}

var stats struct {
	changes    int
	files      int
	badqueries map[string]bool
}

func init() {
	stats.badqueries = make(map[string]bool)
}

func ShowStats() {
	log.Printf("%d files / %d changes / %d bad queries", stats.files, stats.changes, len(stats.badqueries))
}
