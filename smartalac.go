package sa

import (
	"log"
)

var dryrun bool
var debug bool
var skipmb bool
var skipmove bool
var overwrite bool
var finaldir = "/home/music/alac"

func Dryrun(b bool) {
	dryrun = b
}

func Debug(b bool) {
	debug = b
}

func Overwrite(b bool) {
	overwrite = b
}

func Finaldir(f string) {
	finaldir = f
}

func SkipMB(b bool) {
	skipmb = b
}

func SkipMove(b bool) {
	skipmove = b
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
