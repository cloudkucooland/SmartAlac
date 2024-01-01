package sa

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/kr/pretty"
)

var base string

// main entry point
func WalkTree(d string) error {
	base = d
	fsys := os.DirFS(d)
	fs.WalkDir(fsys, ".", wdf)
	return nil
}

// https://pkg.go.dev/io/fs#WalkDirFunc
func wdf(p string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		return nil
	}

	// if jpg/png/etc log in covers...
	if !strings.HasSuffix(p, ".m4a") {
		if debug {
			log.Printf("skipping non-m4a file: %s\n", p)
		}
		return nil
	}

	fullpath := path.Join(base, p)
	log.Println(fullpath)
	mp4, err := mp4tag.Open(fullpath)
	if err != nil {
		log.Printf("unable to open mp4 file: %s", err.Error())
		return err
	}
	defer mp4.Close()

	tags, err := mp4.Read()
	if err != nil {
		log.Println("unable to read mp4 metadata: %s", err.Error())
		return err
	}

	if debug {
		log.Printf("%# v\n", pretty.Formatter(tags.Custom))
	}

	// if already tagged with MBIDs
	tid, ok := tags.Custom["MusicBrainz Album Id"]
	if !ok || tid == "" {
		log.Printf("not yet tagged with MBIDs, skipping (will write interface to query mb later): %s\n", p)
		return nil
	}
	if len(tid) != 36 {
		log.Printf("corrupt MBID [%s]: %s\n", tid, p)
		return nil
	}
	stats.files = stats.files + 1

	newtags, changed, err := updateFromMB(tags)
	if changed {
		stats.changes = stats.changes + showDiffs(tags, newtags)
		if !dryrun {
			if err := mp4.Write(newtags, []string{}); err != nil {
				log.Println("saving : %s", err.Error())
				return err
			}
			// rename the file if needed
		} else {
			log.Printf("Would have saved if not in dry-run mode: %s\n", newtags.Title)
		}
	}

	return nil
}

func showDiffs(in, out *mp4tag.MP4Tags) int {
	d := pretty.Diff(in, out)
	for _, v := range d {
		sp := strings.SplitN(v, ":", 2)
		fmt.Printf("%s\t\t\t%s\n", sp[0], sp[1])
	}
	return len(d)
}
