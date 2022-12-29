package sa

import (
	"io/fs"
	"log"
	"os"
	"path"
)

var rename map[string]string
var covers map[string]string
var base string

func init() {
	rename = make(map[string]string)
	covers = make(map[string]string)
}

// phase 1, walk and build list -- OK to query MB at this point, even to update
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
		// log.Printf("%+v\n", d)
		return nil
	}

	fullpath := path.Join(base, p)

	// if jpg/png/etc log in covers...
	// if mp4, parse tags and do the work
	LoadALACTags(fullpath)
	return nil
}

// phase 2 move files if needed.
// run rename at the end of the process to move files to their new places
func Rename() error {
	for k, v := range rename {
		log.Println(k, v)
		// rename k to v
	}
	return nil
}
