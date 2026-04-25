package sa

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/jo-hoe/chromaprint"
	"github.com/kr/pretty"
)

// main entry point
func (c *Curator) WalkTree(d string) error {
	if d == "" {
		log.Printf("directory string empty")
		return nil
	}

	return filepath.WalkDir(d, c.wdf)
}

// https://pkg.go.dev/io/fs#WalkDirFunc
func (c *Curator) wdf(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d == nil || d.IsDir() {
		return nil
	}

	// if jpg/png/etc log in covers...
	if !strings.HasSuffix(p, ".m4a") || strings.HasPrefix(filepath.Base(p), "._") {
		if c.Config.Debug {
			log.Printf("skipping non-m4a file: %s\n", p)
		}
		return nil
	}

	if c.Config.Debug {
		log.Println(p)
	}
	mp4, err := mp4tag.Open(p)
	if err != nil {
		log.Printf("unable to open mp4 file: %s %s", err.Error(), p)
		return nil // err
	}
	defer mp4.Close()
	mp4.UpperCustom(false)

	tags, err := mp4.Read()
	if err != nil {
		log.Printf("unable to read mp4 metadata: %s %s", err.Error(), p)
		return nil // err
	}

	if c.Config.Debug {
		log.Printf("%# v\n", pretty.Formatter(tags.Custom))
	}
	// if already tagged with MBIDs
	tid := tags.Custom["MusicBrainz Album Id"]
	discID := tags.Custom["MusicBrainz Disc Id"]
	toc := tags.Custom["TOC"]

	if tid == "" {
		// Try to extract MBID or DiscID from directory name: /path/to/Album [MBID]/...
		dir := filepath.Base(filepath.Dir(p))

		// Check for [MBID] format
		if idxStart := strings.LastIndex(dir, "["); idxStart != -1 {
			if idxEnd := strings.LastIndex(dir, "]"); idxEnd != -1 && idxEnd > idxStart+1 {
				potentialID := dir[idxStart+1 : idxEnd]
				if len(potentialID) == 36 {
					tid = potentialID
					if c.Config.Debug {
						log.Printf("extracted MBID from directory: %s\n", tid)
					}
				}
			}
		}

		// Or if the directory name itself is a DiscID (28 chars, base64-ish)
		if tid == "" && discID == "" && len(dir) == 28 {
			discID = dir
			if c.Config.Debug {
				log.Printf("found potential DiscID in directory name: %s\n", discID)
			}
		}
	}

	if tid == "" && discID == "" && toc == "" {
		// Final Fallback: AcoustID Fingerprinting
		if c.Config.AcoustIDKey != "" {
			if c.Config.Debug {
				log.Printf("attempting AcoustID lookup for %s\n", p)
			}
			chromaprinter, err := chromaprint.NewBuilder().Build()
			if err == nil {
				fingerprints, err := chromaprinter.CreateFingerprints(p)
				if err == nil && len(fingerprints) > 0 {
					// Use the first fingerprint (AcoustID usually only needs one)
					fp := fingerprints[0]
					encoded := EncodeFingerprint(fp.Fingerprint)
					resolvedID, err := c.acoustIDLookup(encoded, int(fp.DurationInSeconds))
					if err == nil && resolvedID != "" {
						tid = resolvedID
						if c.Config.Debug {
							log.Printf("resolved AcoustID to release %s\n", tid)
						}
					}
				}
			}
		}
	}

	if tid == "" && discID == "" && toc == "" {
		if c.Config.SkipMB {
			log.Printf("not tagged with MBIDs, skipping: %s\n", p)
			return nil
		}

		artist := tags.AlbumArtist
		if artist == "" {
			artist = tags.Artist
		}
		album := tags.Album
		if album == "" {
			log.Printf("not tagged with MBIDs and no album tag, skipping: %s\n", p)
			return nil
		}

		log.Printf("not tagged with MBIDs, searching for %s - %s\n", artist, album)
		results, err := c.SearchMB(artist, album)
		if err != nil {
			log.Printf("search failed: %v\n", err)
			return nil
		}

		if len(results) == 0 {
			log.Printf("no matches found for %s - %s\n", artist, album)
			return nil
		}

		selected, err := c.SelectRelease(results)
		if err != nil || selected == nil {
			return nil
		}

		tid = selected.ID
	}

	if tid == "" && discID == "" && toc == "" {
		log.Printf("no way to resolve metadata for %s\n", p)
		return nil
	}

	c.Stats.mu.Lock()
	c.Stats.Files++
	c.Stats.mu.Unlock()

	renametags := tags

	if !c.Config.SkipMB {
		var newtags *mp4tag.MP4Tags
		var changed bool
		var err error

		if tid != "" {
			newtags, changed, err = c.updateFromMB(tags, tid)
		} else if discID != "" {
			newtags, changed, err = c.updateFromDiscID(tags, discID)
		}

		renametags = newtags
		if err != nil {
			log.Printf("updating: %s\n", err.Error())
			return err
		}
		if changed {
			diffs := showDiffs(tags, newtags)
			c.Stats.mu.Lock()
			c.Stats.Changes += diffs
			c.Stats.mu.Unlock()

			if c.Config.DryRun {
				log.Printf("Would have saved if not in dry-run mode: %s\n", newtags.Title)
				return nil
			}

			if err := mp4.Write(newtags, []string{}); err != nil {
				log.Printf("error while saving: %s\n", err.Error())
				return err
			}
		}
	}

	if !c.Config.SkipMove {
		if err := c.rename(p, renametags); err != nil {
			log.Printf("error while renaming: %s\n", err.Error())
			return err
		}
	}

	return nil
}

func showDiffs(in, out *mp4tag.MP4Tags) int {
	d := pretty.Diff(in, out)
	for _, v := range d {
		sp := strings.SplitN(v, ":", 2)
		fmt.Printf("%s\t%s\n", sp[0], sp[1])
	}
	return len(d)
}

func (c *Curator) rename(fullpath string, tags *mp4tag.MP4Tags) error {
	if tags.AlbumArtist == "" {
		return errors.New("albumArtist not set, not moving")
	}
	if tags.AlbumSort == "" {
		return errors.New("AlbumSort not set, not moving")
	}

	aa := tags.AlbumArtistSort
	if aa == "" {
		aa = tags.AlbumArtist
	}
	aa = sanitize(aa)

	artistdir := filepath.Join(c.Config.FinalDir, aa)
	if err := os.MkdirAll(artistdir, 0755); err != nil {
		return err
	}

	album := sanitize(tags.AlbumSort)
	albumdir := filepath.Join(artistdir, album)
	if err := os.MkdirAll(albumdir, 0755); err != nil {
		return err
	}

	cleantitle := sanitize(tags.Title)
	if len(cleantitle) > 100 {
		cleantitle = cleantitle[0:100]
	}

	filename := fmt.Sprintf("%d-%02d %s.m4a", tags.DiscNumber, tags.TrackNumber, cleantitle)
	finalpath := filepath.Join(albumdir, filename)

	if filepath.Clean(finalpath) == filepath.Clean(fullpath) {
		if c.Config.Debug {
			log.Printf("no need to move: %s\n", fullpath)
		}
		return nil
	}

	_, err := os.Stat(finalpath)
	if err == nil && !c.Config.Overwrite {
		log.Printf("file already exists, not overwriting: %s (from %s)\n", finalpath, fullpath)
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Printf("moving %s to %s\n", fullpath, finalpath)
	if c.Config.DryRun {
		return nil
	}

	if err := os.Rename(fullpath, finalpath); err != nil {
		return c.move(fullpath, finalpath)
	}

	return nil
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "?", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, ">", "_")
	s = strings.ReplaceAll(s, "<", "_")
	s = strings.ReplaceAll(s, "|", "_")
	s = strings.ReplaceAll(s, "*", "_")
	s = strings.ReplaceAll(s, "\"", "_")
	s = strings.ReplaceAll(s, "'", "_")
	return s
}

// cross-mountpoint safe move
func (c *Curator) move(oldpath, newpath string) error {
	source, err := os.Open(oldpath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(newpath)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	source.Close() // close before remove
	return os.Remove(oldpath)
}
