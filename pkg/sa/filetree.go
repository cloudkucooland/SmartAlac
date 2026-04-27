package sa

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/jo-hoe/chromaprint"
	"github.com/kr/pretty"
)

// main entry point
func (c *Curator) WalkTree(ctx context.Context, d string) error {
	if d == "" {
		log.Printf("directory string empty")
		return nil
	}
	c.ctx = ctx

	return filepath.WalkDir(d, c.wdf)
}

// WalkDirFunc
func (c *Curator) wdf(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Check if file still exists (it might have been moved by a previous step in the walk)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}

	// Check context for graceful shutdown before starting work on a new file
	if c.ctx != nil {
		select {
		case <-c.ctx.Done():
			log.Printf("stopping curation: %v\n", c.ctx.Err())
			return c.ctx.Err()
		default:
		}
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
		return nil // continue
	}
	defer mp4.Close()
	mp4.UpperCustom(false)

	inTags, err := mp4.Read()
	if err != nil {
		log.Printf("unable to read mp4 metadata: %s %s", err.Error(), p)
		return nil // continue
	}

	if c.Config.Debug {
		log.Printf("%# v\n", pretty.Formatter(inTags.Custom))
	}

	// Working state
	tags := inTags
	tid := tags.Custom["MusicBrainz Album Id"]
	discID := tags.Custom["MusicBrainz Disc Id"]
	toc := tags.Custom["TOC"]
	did := tags.Custom["Discogs Release Id"]

	// 1. Try to find IDs from directory name if missing
	if tid == "" && discID == "" && did == "" {
		dir := filepath.Base(filepath.Dir(p))
		if idxStart := strings.LastIndex(dir, "["); idxStart != -1 {
			if idxEnd := strings.LastIndex(dir, "]"); idxEnd != -1 && idxEnd > idxStart+1 {
				potentialID := dir[idxStart+1 : idxEnd]
				if isMBID(potentialID) {
					tid = potentialID
				} else if strings.HasPrefix(potentialID, "discogs:") {
					did = strings.TrimPrefix(potentialID, "discogs:")
				}
			}
		}
	}

	// 2. AcoustID Fingerprinting Fallback
	if tid == "" && discID == "" && did == "" && c.Config.AcoustIDKey != "" {
		var fpStr string
		var dur int
		var resolvedAID string

		info, _ := d.Info()
		if c.Cache != nil {
			fpStr, dur, resolvedAID, _ = c.Cache.GetFingerprint(p, info.ModTime(), info.Size())
		}

		if fpStr == "" {
			chromaprinter, err := chromaprint.NewBuilder().WithPathToChromaprint(c.Config.FpcalcPath).Build()
			if err == nil {
				fingerprints, err := chromaprinter.CreateFingerprints(p)
				if err == nil && len(fingerprints) > 0 {
					f := fingerprints[0]
					fpStr = EncodeFingerprint(f.Fingerprint)
					dur = int(f.DurationInSeconds)
					resolvedAID, _ = c.AcoustIDLookup(c.ctx, fpStr, dur)
					if c.Cache != nil {
						c.Cache.SaveFingerprint(p, info.ModTime(), info.Size(), fpStr, dur, resolvedAID)
					}
				}
			}
		}

		if resolvedAID != "" {
			tid = resolvedAID
		}
	}

	// 3. Search MB Fallback
	if tid == "" && discID == "" && did == "" && !c.Config.SkipMB {
		artist := tags.AlbumArtist
		if artist == "" { artist = tags.Artist }
		album := tags.Album
		if artist != "" && album != "" {
			results, err := c.SearchMB(c.ctx, artist, album)
			if err == nil && len(results) > 0 {
				selected, _ := c.SelectRelease(c.ctx, results)
				if selected != nil {
					tid = selected.ID
				}
			}
		}
	}

	// 4. Search Discogs Fallback
	if tid == "" && discID == "" && did == "" && c.dgClient != nil {
		artist := tags.AlbumArtist
		if artist == "" { artist = tags.Artist }
		album := tags.Album
		if artist != "" && album != "" {
			results, err := c.SearchDiscogs(c.ctx, artist, album)
			if err == nil && len(results) > 0 {
				did = strconv.Itoa(results[0].ID)
			}
		}
	}

	// 5. Update Tags from Resolved IDs
	var finalTags *mp4tag.MP4Tags
	var changed bool

	if tid != "" {
		finalTags, changed, err = c.UpdateFromMB(c.ctx, inTags, tid)
	} else if discID != "" {
		finalTags, changed, err = c.UpdateFromDiscID(c.ctx, inTags, discID)
	} else if did != "" {
		didInt, _ := strconv.Atoi(did)
		finalTags, changed, err = c.UpdateFromDiscogs(c.ctx, inTags, didInt)
	} else if toc != "" {
		// Just fresh data from MB if TOC exists
		finalTags, changed, err = c.UpdateFromMB(c.ctx, inTags, "")
	}

	if err != nil {
		log.Printf("metadata update failed for %s: %v", p, err)
		return nil
	}

	if finalTags == nil {
		log.Printf("no way to resolve metadata for %s", p)
		return nil
	}

	c.Stats.mu.Lock()
	c.Stats.Files++
	c.Stats.mu.Unlock()

	if changed {
		diffs := showDiffs(inTags, finalTags)
		c.Stats.mu.Lock()
		c.Stats.Changes += diffs
		c.Stats.mu.Unlock()

		if !c.Config.DryRun {
			customKeys := make([]string, 0, len(finalTags.Custom))
			for k := range finalTags.Custom {
				customKeys = append(customKeys, k)
			}
			if err := mp4.Write(finalTags, customKeys); err != nil {
				log.Printf("error saving tags: %v", err)
				return err
			}
		}
	}

	if !c.Config.SkipMove {
		if err := c.rename(p, finalTags); err != nil {
			log.Printf("error renaming: %v", err)
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
		return errors.New("AlbumArtist not set, not moving")
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

	srcStat, err := os.Stat(fullpath)
	if err != nil {
		return err
	}

	destStat, err := os.Stat(finalpath)
	if err == nil {
		if os.SameFile(srcStat, destStat) {
			if c.Config.Debug {
				log.Printf("no need to move (same file): %s\n", fullpath)
			}
			return nil
		}
		if !c.Config.Overwrite {
			log.Printf("file already exists, not overwriting: %s (from %s)\n", finalpath, fullpath)
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	fmt.Printf("moving %s to %s\n", fullpath, finalpath)
	if c.Config.DryRun {
		return nil
	}

	if err := os.Rename(fullpath, finalpath); err != nil {
		if err := c.move(fullpath, finalpath); err != nil {
			log.Printf("error moving file: %v", err)
			return nil // Log and continue walk
		}
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
