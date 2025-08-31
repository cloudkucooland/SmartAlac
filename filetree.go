package sa

import (
	"errors"
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
	if d == "" {
		log.Printf("directory string empty")
		return nil
	}

	base = d
	fsys := os.DirFS(d)
	fs.WalkDir(fsys, ".", wdf)
	return nil
}

// https://pkg.go.dev/io/fs#WalkDirFunc
func wdf(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d == nil || d.IsDir() {
		return nil
	}

	// if jpg/png/etc log in covers...
	if !strings.HasSuffix(p, ".m4a") || strings.HasPrefix(p, "._") {
		if debug {
			log.Printf("skipping non-m4a file: %s\n", p)
		}
		return nil
	}

	fullpath := path.Join(base, p)
	if debug {
		log.Println(fullpath)
	}
	mp4, err := mp4tag.Open(fullpath)
	if err != nil {
		log.Printf("unable to open mp4 file: %s %s", err.Error(), fullpath)
		return nil // err
	}
	defer mp4.Close()
	mp4.UpperCustom(false)

	tags, err := mp4.Read()
	if err != nil {
		log.Printf("unable to read mp4 metadata: %s %s", err.Error(), fullpath)
		return nil // err
	}

	if debug {
		log.Printf("%# v\n", pretty.Formatter(tags.Custom))
	}

	// if already tagged with MBIDs
	tid, ok := tags.Custom["MusicBrainz Album Id"]
	if !ok || tid == "" {
		log.Printf("not tagged with MBIDs, skipping (will write interface to query mb later): %s\n", p)
		return nil
	}
	if len(tid) != 36 {
		log.Printf("corrupt MBID [%s]: %s\n", tid, p)
		return nil
	}
	stats.files = stats.files + 1

	newtags, changed, err := updateFromMB(tags)
	if err != nil {
		log.Printf("updating: %s\n", err.Error())
		return err
	}
	if changed {
		stats.changes = stats.changes + showDiffs(tags, newtags)
		if dryrun {
			log.Printf("Would have saved if not in dry-run mode: %s\n", newtags.Title)
			return nil
		}

		if err := mp4.Write(newtags, []string{}); err != nil {
			log.Printf("error while saving: %s\n", err.Error())
			return err
		}

		if err := rename(fullpath, newtags); err != nil {
			log.Printf("error while renaming: %s\n", err.Error())
			return err
		}
	}

	return nil
}

// it is an anti-pattern to overload a function in this way: one function to show, another to count please
func showDiffs(in, out *mp4tag.MP4Tags) int {
	d := pretty.Diff(in, out)
	/* for _, v := range d {
		sp := strings.SplitN(v, ":", 2)
		fmt.Printf("%s\t%s\n", sp[0], sp[1])
	} */
	return len(d)
}

func rename(fullpath string, tags *mp4tag.MP4Tags) error {
	if tags.AlbumArtist == "" {
		return errors.New("Artist not set, not moving")
	}
	if tags.AlbumSort == "" {
		return errors.New("AlbumSort not set, not moving")
	}

	aa := tags.AlbumArtistSort
	if aa == "" {
		aa = tags.AlbumArtist
	}
	aa = strings.ReplaceAll(aa, "/", "_")
	aa = strings.ReplaceAll(aa, "?", "_")
	aa = strings.ReplaceAll(aa, ":", "_")
	aa = strings.ReplaceAll(aa, ">", "_")
	aa = strings.ReplaceAll(aa, "\"", "_")
	aa = strings.ReplaceAll(aa, "'", "_")

	artistdir := path.Join(finaldir, aa)
	artistdirstat, err := os.Stat(artistdir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("making artistdir: %s\n", artistdir)
			if err := os.Mkdir(artistdir, 0755); err != nil {
				return err
			}
			artistdirstat, _ = os.Stat(artistdir)
		} else {
			return err
		}
	}
	if !artistdirstat.IsDir() {
		return fmt.Errorf("artist directory is not directory... %s", artistdir)
	}

	album := tags.AlbumSort
	album = strings.ReplaceAll(album, "/", "_")
	album = strings.ReplaceAll(album, "?", "_")
	album = strings.ReplaceAll(album, ":", "_")
	album = strings.ReplaceAll(album, ">", "_")
	album = strings.ReplaceAll(album, "\"", "_")
	album = strings.ReplaceAll(album, "'", "_")

	albumdir := path.Join(artistdir, album)
	albumdirstat, err := os.Stat(albumdir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("making albumdir: %s\n", albumdir)
			if err := os.Mkdir(albumdir, 0755); err != nil {
				return err
			}
			albumdirstat, _ = os.Stat(albumdir)
		} else {
			return err
		}
	}
	if !albumdirstat.IsDir() {
		return fmt.Errorf("album  directory is not directory... %s", albumdir)
	}

	cleantitle := strings.ReplaceAll(tags.Title, "/", "_")
	cleantitle = strings.ReplaceAll(cleantitle, "?", "_")
	cleantitle = strings.ReplaceAll(cleantitle, ":", "_")
	cleantitle = strings.ReplaceAll(cleantitle, ">", "_")
	cleantitle = strings.ReplaceAll(cleantitle, "\"", "_")
	cleantitle = strings.ReplaceAll(cleantitle, "'", "_")

	if len(cleantitle) > 100 {
		cleantitle = cleantitle[0:100]
	}

	filename := fmt.Sprintf("%d-%02d %s.m4a", tags.DiscNumber, tags.TrackNumber, cleantitle)
	finalpath := path.Join(albumdir, filename)
	if finalpath == fullpath {
		if debug {
			log.Printf("no need to move: %s\n", fullpath)
		}
		return nil
	}

	_, err = os.Stat(finalpath)
	if err == nil {
		log.Printf("file already exists, not overwriting: %s (from %s)\n", finalpath, fullpath)
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Printf("moving %s to %s\n", fullpath, finalpath)
	if err := os.Rename(fullpath, finalpath); err != nil {
		return err
	}

	return nil
}
