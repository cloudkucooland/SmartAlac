package sa

import (
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
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
		log.Printf("skipping non-m4a file: %s\n", p)
		return nil
	}

	fullpath := path.Join(base, p)
	log.Println(fullpath)
	mp4, err := mp4tag.Open(fullpath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tags, err := mp4.Read()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// if already tagged with MBIDs
	if tid, ok := tags.Custom["MusicBrainz Release Track Id"]; !ok || tid == "" {
		log.Printf("not yet tagged with MBIDs, skiipping (for now): %s\n", p)
		return nil
	}

	newtags, changed, err := updateFromMB(tags)
	if changed {
		/* if err := mp4.Write(newtags, []string{}); err != nil {
			log.Println(err.Error())
			return err
		} */
		log.Printf("Would have saved if it were enabled: %s\n", newtags.Title)
		// rename the file if needed
	}

	return nil
}
