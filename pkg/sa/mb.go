package sa

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
	"github.com/kr/pretty"
)

type AcoustIDResponse struct {
	Status  string `json:"status"`
	Results []struct {
		ID         string  `json:"id"`
		Score      float64 `json:"score"`
		Recordings []struct {
			ID            string `json:"id"`
			Releasegroups []struct {
				ID       string `json:"id"`
				Releases []struct {
					ID string `json:"id"`
				} `json:"releases"`
			} `json:"releasegroups"`
		} `json:"recordings"`
	} `json:"results"`
}

func isMBID(id string) bool {
	if len(id) != 36 {
		return false
	}
	// Basic UUID check: 8-4-4-4-12
	return strings.Count(id, "-") == 4
}

func (c *Curator) UpdateFromMB(ctx context.Context, in *mp4tag.MP4Tags, overrideID string) (*mp4tag.MP4Tags, bool, error) {
	if c.mb5query == nil {
		return in, false, fmt.Errorf("mb5 query not initialized")
	}

	select {
	case <-ctx.Done():
		return in, false, ctx.Err()
	default:
	}

	if c.Config.Debug {
		temp := *in
		temp.Pictures = nil
		log.Printf("%# v", pretty.Formatter(temp))
	}

	releaseid := overrideID
	if releaseid == "" {
		var ok bool
		releaseid, ok = in.Custom["MusicBrainz Album Id"]
		if !ok {
			log.Println("no release ID, skipping")
			return in, false, nil
		}
	}

	// Validate releaseid. If it's malformed and was in the tags, we want to remove it.
	if !isMBID(releaseid) {
		if releaseid != "" && releaseid != "0" {
			log.Printf("malformed MusicBrainz ID found and removed: %q", releaseid)
			out := *in
			if out.Custom != nil {
				newCustom := make(map[string]string)
				for k, v := range out.Custom {
					newCustom[k] = v
				}
				delete(newCustom, "MusicBrainz Album Id")
				// Also scrub other potentially malformed MBIDs
				for k, v := range newCustom {
					if strings.HasPrefix(k, "MusicBrainz") && !isMBID(v) {
						delete(newCustom, k)
					}
				}
				out.Custom = newCustom
				return &out, true, nil
			}
		}
		return in, false, nil
	}

	recordingID := in.Custom["MusicBrainz Track Id"]

	// Track and Disc numbers are required for matching within a release
	if in.TrackNumber < 1 {
		log.Printf("skipping track with number %d", in.TrackNumber)
		return in, false, nil
	}
	if in.DiscNumber < 1 {
		// Default to disc 1 if not set
		in.DiscNumber = 1
	}

	c.Stats.mu.Lock()
	if v, ok := c.Stats.BadQueries[releaseid]; ok && v {
		c.Stats.mu.Unlock()
		return in, false, fmt.Errorf("%s failed once already, skipping", releaseid)
	}
	c.Stats.mu.Unlock()

	var metadata mb5.Metadata
	var retryCount int

	// TTL Cache check: If we've looked this up recently, don't hit the API again
	if c.Cache != nil {
		_, lastChecked, _, err := c.Cache.GetRelease(releaseid)
		if err == nil && time.Since(lastChecked) < 30*24*time.Hour {
			if c.Config.Debug {
				log.Printf("skipping recently checked release %s (cached %v ago)", releaseid, time.Since(lastChecked))
			}
			return in, false, nil
		}
	}

	c.rl.Take()
	for {
		select {
		case <-ctx.Done():
			return in, false, ctx.Err()
		default:
		}

		// Prepare query parameters for explicit "inc" request
		p1 := []byte("inc\x00")
		v1 := []byte("artists labels recordings release-groups url-rels artist-credits work-rels artist-rels work-level-rels\x00")

		params := [1]*byte{&p1[0]}
		values := [1]*byte{&v1[0]}

		// Query libmusicbrainz5
		metadata = mb5.QueryQuery(c.mb5query, "release", releaseid, "", 1, unsafe.Pointer(&params[0]), unsafe.Pointer(&values[0]))
		if metadata != nil {
			// Mark as checked in Cache
			if c.Cache != nil {
				c.Cache.SaveRelease(releaseid, nil, "")
			}
			break
		}
		lastCode := mb5.QueryGetLasthttpcode(c.mb5query)
		if lastCode == 503 && retryCount < 3 {
			retryCount++
			log.Printf("MusicBrainz returned 503, retrying in 5 seconds (attempt %d/3)...", retryCount)

			// Context-aware sleep
			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return in, false, ctx.Err()
			}
			continue
		}

		c.Stats.mu.Lock()
		c.Stats.BadQueries[releaseid] = true
		c.Stats.mu.Unlock()

		var errbuf [256]byte
		mb5.QueryGetLasterrormessage(c.mb5query, &errbuf[0], 256)
		cErr := strings.Trim(string(errbuf[:]), "\x00")
		log.Printf("query to MusicBrainz failed for %s (HTTP %d): %s", releaseid, lastCode, cErr)
		return in, false, nil
	}

	defer mb5.MetadataDelete(metadata)
	release := mb5.MetadataGetRelease(metadata)

	if release == nil {
		return in, false, fmt.Errorf("no release in metadata for %s", releaseid)
	}

	// Start with a copy of input tags to preserve existing data
	out := *in
	if out.Custom == nil {
		out.Custom = make(map[string]string)
	} else {
		// Deep copy the map to avoid side effects
		newCustom := make(map[string]string)
		for k, v := range out.Custom {
			newCustom[k] = v
		}
		out.Custom = newCustom
	}

	out.Album = mb5.String(mb5.ReleaseGetTitle, unsafe.Pointer(release))
	out.AlbumSort = out.Album + " [" + releaseid + "]"

	if disambig := mb5.String(mb5.ReleaseGetDisambiguation, unsafe.Pointer(release)); disambig != "" {
		out.Comment = disambig
	}

	ac := mb5.ReleaseGetArtistcredit(release)
	if aa := c.fmtArtistCreditMB5(ac); aa != "" {
		out.AlbumArtist = aa
	}
	if asa := c.fmtArtistCreditSortMB5(ac); asa != "" && asa != out.AlbumArtist {
		out.AlbumArtistSort = asa
	}

	mediumList := mb5.ReleaseGetMediumlist(release)
	out.DiscTotal = int16(mb5.MediumListSize(mediumList))

	var foundTrack mb5.Track
	var foundMedium mb5.Medium
	for i := 0; i < int(out.DiscTotal); i++ {
		m := mb5.MediumListItem(mediumList, i)
		if mb5.MediumGetPosition(m) == int(in.DiscNumber) {
			foundMedium = m
			tl := mb5.MediumGetTracklist(m)
			out.TrackTotal = int16(mb5.TrackListSize(tl))
			for j := 0; j < int(out.TrackTotal); j++ {
				t := mb5.TrackListItem(tl, j)
				if mb5.TrackGetPosition(t) == int(in.TrackNumber) {
					foundTrack = t
					break
				}
			}
			break
		}
	}

	if foundTrack != nil {
		recording := mb5.TrackGetRecording(foundTrack)
		if recording != nil {
			if recordingID == "" {
				var buf [37]byte
				mb5.RecordingGetID(unsafe.Pointer(recording), (*byte)(unsafe.Pointer(&buf[0])), 37)
				recordingID = strings.Trim(string(buf[:]), "\x00")
			}
			tac := mb5.RecordingGetArtistcredit(recording)
			if tac == nil {
				tac = mb5.TrackGetArtistcredit(foundTrack)
			}
			if artist := c.fmtArtistCreditMB5(tac); artist != "" {
				out.Artist = artist
			}

			title := mb5.String(mb5.RecordingGetTitle, unsafe.Pointer(recording))
			if title == "" {
				title = mb5.String(mb5.TrackGetTitle, unsafe.Pointer(foundTrack))
			}
			if title != "" {
				out.Title = title
			}

			if alist := c.fmtArtistListMB5(tac); alist != "" {
				out.Custom["ARTISTS"] = alist
			}
			if aids := c.joinArtistIDsMB5(tac); aids != "" {
				out.Custom["MusicBrainz Artist Id"] = aids
			}

			c.processRelations(recording, out.Custom, &out)
			if isrcs := c.fmtISRCsMB5(recording); isrcs != "" {
				out.Custom["ISRC"] = isrcs
			}
		}
	}

	rg := mb5.ReleaseGetReleasegroup(release)
	if rg != nil {
		origDate := mb5.String(mb5.ReleasegroupGetFirstreleasedate, unsafe.Pointer(rg))
		if len(origDate) >= 4 {
			out.Custom["ORIGINALDATE"] = origDate
			out.Custom["ORIGINALYEAR"] = origDate[:4]
			fmt.Sscanf(origDate[:4], "%d", &out.Year)
		}
		if rgid := mb5.String(mb5.ReleasegroupGetID, unsafe.Pointer(rg)); rgid != "" {
			out.Custom["MusicBrainz Release Group Id"] = rgid
		}
	}

	if date := mb5.String(mb5.ReleaseGetDate, unsafe.Pointer(release)); date != "" {
		out.Date = date
	}

	if asin := mb5.String(mb5.ReleaseGetAsin, unsafe.Pointer(release)); asin != "" {
		out.Custom["ASIN"] = asin
	}
	if barcode := mb5.String(mb5.ReleaseGetBarcode, unsafe.Pointer(release)); barcode != "" {
		out.Custom["BARCODE"] = barcode
	}

	liList := mb5.ReleaseGetLabelinfolist(release)
	if liList != nil {
		if label := c.fmtLabelsMB5(liList); label != "" {
			out.Custom["LABEL"] = label
		}
		if cat := c.fmtCatalogNumbersMB5(liList); cat != "" {
			out.Custom["CATALOGNUMBER"] = cat
		}
	}

	countryCode := mb5.String(mb5.ReleaseGetCountry, unsafe.Pointer(release))
	if countryCode != "" {
		if countryName := resolveCountry(countryCode); countryName != "" {
			out.Custom["Country"] = countryName
		}
		out.Custom["MusicBrainz Album Release Country"] = countryCode
	}

	if foundMedium != nil {
		if format := mb5.String(mb5.MediumGetFormat, unsafe.Pointer(foundMedium)); format != "" {
			out.Custom["MEDIA"] = mediumFormat(format)
		}
	}

	if aids := c.joinArtistIDsMB5(ac); aids != "" {
		out.Custom["MusicBrainz Album Artist Id"] = aids
	}
	out.Custom["MusicBrainz Album Id"] = releaseid
	out.Custom["MusicBrainz Track Id"] = recordingID

	// Fallback: If AlbumArtist is still empty, use Artist
	if out.AlbumArtist == "" && out.Artist != "" {
		out.AlbumArtist = out.Artist
	}
	if out.AlbumArtistSort == "" && out.AlbumArtist != "" {
		out.AlbumArtistSort = out.AlbumArtist
	}

	if c.Config.Debug {
		temp := out
		temp.Pictures = nil
		log.Printf("AFTER MB UPDATE: %# v", pretty.Formatter(temp))
	}

	changed := !tagsEquivalent(in, &out)

	return &out, changed, nil
}

func normalizeString(s string) string {
	// Normalize common MusicBrainz/Smart punctuation to standard ASCII
	r := strings.NewReplacer(
		"’", "'",
		"‘", "'",
		"´", "'",
		"”", "\"",
		"“", "\"",
		"‐", "-", // hyphen
		"−", "-", // minus
		"–", "-", // en-dash
		"—", "-", // em-dash
		"…", "...",
		"\u00a0", " ", // non-breaking space
	)
	return strings.ToLower(strings.TrimSpace(r.Replace(s)))
}

func tagsEquivalent(in, out *mp4tag.MP4Tags) bool {
	// Deep compare custom maps first (with normalization for values)
	if len(in.Custom) != len(out.Custom) {
		return false
	}
	for k, v := range in.Custom {
		if normalizeString(out.Custom[k]) != normalizeString(v) {
			return false
		}
	}

	// Compare Date with specificity tolerance (e.g., "1995" matches "1995-06-13")
	if in.Date != out.Date {
		if in.Date != "" && out.Date != "" {
			if !strings.HasPrefix(out.Date, in.Date) {
				return false
			}
		} else {
			return false
		}
	}

	// Compare primary strings with normalization
	if normalizeString(in.Title) != normalizeString(out.Title) {
		return false
	}
	if normalizeString(in.Album) != normalizeString(out.Album) {
		return false
	}
	if normalizeString(in.Artist) != normalizeString(out.Artist) {
		return false
	}
	if normalizeString(in.AlbumArtist) != normalizeString(out.AlbumArtist) {
		return false
	}

	// Create copies for field comparison (excluding fields we already checked)
	inCopy := *in
	outCopy := *out
	inCopy.Pictures = nil
	outCopy.Pictures = nil
	inCopy.Custom = nil
	outCopy.Custom = nil
	inCopy.Date = ""
	outCopy.Date = ""
	inCopy.Title = ""
	outCopy.Title = ""
	inCopy.Album = ""
	outCopy.Album = ""
	inCopy.Artist = ""
	outCopy.Artist = ""
	inCopy.AlbumArtist = ""
	outCopy.AlbumArtist = ""

	return reflect.DeepEqual(inCopy, outCopy)
}

func (c *Curator) processRelations(recording mb5.Recording, custom map[string]string, out *mp4tag.MP4Tags) {
	rll := mb5.RecordingGetRelationlistlist(recording)
	if rll == nil {
		return
	}

	roles := map[string]string{
		"composer":  "COMPOSER",
		"lyricist":  "LYRICIST",
		"producer":  "PRODUCER",
		"engineer":  "ENGINEER",
		"mixer":     "MIXER",
		"remixer":   "REMIXER",
		"writer":    "WRITER",
		"arranger":  "ARRANGER",
		"conductor": "CONDUCTOR",
	}

	credits := make(map[string][]string)

	rllSize := mb5.RelationlistListSize(rll)
	for i := 0; i < rllSize; i++ {
		rl := mb5.RelationlistListItem(rll, i)
		rlSize := mb5.RelationListSize(rl)
		for j := 0; j < rlSize; j++ {
			rel := mb5.RelationListItem(rl, j)
			relType := mb5.String(mb5.RelationGetType, unsafe.Pointer(rel))

			if role, ok := roles[relType]; ok {
				artist := mb5.RelationGetArtist(rel)
				if artist != nil {
					name := mb5.String(mb5.ArtistGetName, unsafe.Pointer(artist))
					if name != "" {
						credits[role] = append(credits[role], name)
					}
				}
			}

			if relType == "performance" {
				work := mb5.RelationGetWork(rel)
				if work != nil {
					if wid := mb5.String(mb5.WorkGetID, unsafe.Pointer(work)); wid != "" {
						custom["MusicBrainz Work Id"] = wid
					}
					if wtitle := mb5.String(mb5.WorkGetTitle, unsafe.Pointer(work)); wtitle != "" {
						custom["WORK"] = wtitle
					}
				}
			}
		}
	}

	for role, names := range credits {
		val := strings.Join(names, ", ")
		if val != "" {
			custom[role] = val
			if role == "COMPOSER" {
				out.Composer = val
			}
			if role == "CONDUCTOR" {
				out.Conductor = val
			}
		}
	}
}

func (c *Curator) fmtISRCsMB5(recording mb5.Recording) string {
	isrcList := mb5.RecordingGetISRCList(recording)
	if isrcList == nil {
		return ""
	}
	size := mb5.ISRCListSize(isrcList)
	var isrcs []string
	for i := 0; i < size; i++ {
		isrc := mb5.ISRCListItem(isrcList, i)
		id := mb5.String(mb5.ISRCGetID, unsafe.Pointer(isrc))
		if id != "" {
			isrcs = append(isrcs, id)
		}
	}
	return strings.Join(isrcs, ", ")
}

func (c *Curator) fmtArtistCreditMB5(ac mb5.ArtistCredit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5.ArtistcreditGetNamecreditlist(ac)
	count := mb5.NamecreditListSize(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5.NamecreditListItem(ncl, i)
		name := mb5.String(mb5.NamecreditGetName, unsafe.Pointer(nc))
		if name == "" {
			artist := mb5.NamecreditGetArtist(nc)
			name = mb5.String(mb5.ArtistGetName, unsafe.Pointer(artist))
		}
		join := mb5.String(mb5.NamecreditGetJoinphrase, unsafe.Pointer(nc))
		s += name + join
	}
	return s
}

func (c *Curator) fmtArtistCreditSortMB5(ac mb5.ArtistCredit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5.ArtistcreditGetNamecreditlist(ac)
	count := mb5.NamecreditListSize(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5.NamecreditListItem(ncl, i)
		artist := mb5.NamecreditGetArtist(nc)
		name := mb5.String(mb5.ArtistGetSortname, unsafe.Pointer(artist))
		join := mb5.String(mb5.NamecreditGetJoinphrase, unsafe.Pointer(nc))
		s += name + join
	}
	return s
}

func (c *Curator) fmtArtistListMB5(ac mb5.ArtistCredit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5.ArtistcreditGetNamecreditlist(ac)
	count := mb5.NamecreditListSize(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5.NamecreditListItem(ncl, i)
		artist := mb5.NamecreditGetArtist(nc)
		name := mb5.String(mb5.ArtistGetName, unsafe.Pointer(artist))
		if i > 0 {
			s += ", "
		}
		s += name
	}
	return s
}

func (c *Curator) joinArtistIDsMB5(ac mb5.ArtistCredit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5.ArtistcreditGetNamecreditlist(ac)
	count := mb5.NamecreditListSize(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5.NamecreditListItem(ncl, i)
		artist := mb5.NamecreditGetArtist(nc)
		id := mb5.String(mb5.ArtistGetID, unsafe.Pointer(artist))
		if i > 0 {
			s += ","
		}
		s += id
	}
	return s
}

func (c *Curator) fmtLabelsMB5(liList mb5.LabelInfoList) string {
	count := mb5.LabelinfoListSize(liList)
	var labels []string
	seen := make(map[string]bool)
	for i := 0; i < count; i++ {
		li := mb5.LabelinfoListItem(liList, i)
		label := mb5.LabelinfoGetLabel(li)
		if label != nil {
			name := mb5.String(mb5.LabelGetName, unsafe.Pointer(label))
			if name != "" && !seen[name] {
				labels = append(labels, name)
				seen[name] = true
			}
		}
	}
	return strings.Join(labels, "; ")
}

func (c *Curator) fmtCatalogNumbersMB5(liList mb5.LabelInfoList) string {
	count := mb5.LabelinfoListSize(liList)
	var cats []string
	seen := make(map[string]bool)
	for i := 0; i < count; i++ {
		li := mb5.LabelinfoListItem(liList, i)
		cat := mb5.String(mb5.LabelinfoGetCatalognumber, unsafe.Pointer(li))
		if cat != "" && !seen[cat] {
			cats = append(cats, cat)
			seen[cat] = true
		}
	}
	return strings.Join(cats, "; ")
}

func (c *Curator) UpdateFromDiscID(ctx context.Context, in *mp4tag.MP4Tags, discid string) (*mp4tag.MP4Tags, bool, error) {
	if c.mb5query == nil {
		return in, false, fmt.Errorf("mb5 query not initialized")
	}

	select {
	case <-ctx.Done():
		return in, false, ctx.Err()
	default:
	}

	c.rl.Take()

	var metadata mb5.Metadata
	var retryCount int
	for {
		select {
		case <-ctx.Done():
			return in, false, ctx.Err()
		default:
		}

		metadata = mb5.QueryQuery(c.mb5query, "discid", discid, "", 0, nil, nil)
		if metadata != nil {
			break
		}

		lastCode := mb5.QueryGetLasthttpcode(c.mb5query)
		if lastCode == 503 && retryCount < 3 {
			retryCount++
			log.Printf("MusicBrainz returned 503 (discid), retrying in 5 seconds (attempt %d/3)...", retryCount)
			
			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return in, false, ctx.Err()
			}
			continue
		}

		if lastCode == 404 {
			return in, false, nil
		}

		return in, false, fmt.Errorf("discid lookup failed for %s (HTTP %d)", discid, lastCode)
	}
	defer mb5.MetadataDelete(metadata)

	disc := mb5.MetadataGetDisc(metadata)
	if disc == nil {
		return in, false, nil
	}

	rl := mb5.DiscGetReleaselist(disc)
	if rl == nil || mb5.ReleaseListSize(rl) == 0 {
		return in, false, nil
	}

	// For now, we take the first release and use its ID
	rel := mb5.ReleaseListItem(rl, 0)
	var relBuf [37]byte
	mb5.ReleaseGetID(unsafe.Pointer(rel), (*byte)(unsafe.Pointer(&relBuf[0])), 37)
	releaseID := strings.Trim(string(relBuf[:]), "\x00")

	if c.Config.Debug {
		log.Printf("resolved discid %s to release %s\n", discid, releaseID)
	}

	// Now call the standard update with this release ID
	out, changed, err := c.UpdateFromMB(ctx, in, releaseID)
	if err == nil && out != nil {
		if out.Custom == nil {
			out.Custom = make(map[string]string)
		}
		out.Custom["MusicBrainz Disc Id"] = discid
	}
	return out, changed, err
}

func (c *Curator) AcoustIDLookup(ctx context.Context, fingerprint string, duration int) (string, error) {
	if c.Config.AcoustIDKey == "" {
		return "", fmt.Errorf("AcoustID key not set")
	}

	url := fmt.Sprintf("https://api.acoustid.org/v2/lookup?client=%s&meta=releasegroups+releases&duration=%d&fingerprint=%s",
		c.Config.AcoustIDKey, duration, fingerprint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AcoustID API returned HTTP %d", resp.StatusCode)
	}

	var air AcoustIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&air); err != nil {
		return "", err
	}

	if air.Status != "ok" {
		return "", fmt.Errorf("AcoustID API status: %s", air.Status)
	}

	// Find the first release ID with a good score
	for _, res := range air.Results {
		if res.Score < 0.8 {
			continue
		}
		for _, rec := range res.Recordings {
			for _, rg := range rec.Releasegroups {
				for _, rel := range rg.Releases {
					if rel.ID != "" {
						return rel.ID, nil
					}
				}
			}
		}
	}

	return "", nil
}

// TagDirectory tags all .m4a files in a directory with the given MusicBrainz release ID.
func (c *Curator) TagDirectory(ctx context.Context, dir, releaseID string, p interface{}) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !strings.HasSuffix(f.Name(), ".m4a") {
			continue
		}

		path := filepath.Join(dir, f.Name())
		mp4, err := mp4tag.Open(path)
		if err != nil {
			log.Printf("error opening %s: %v", path, err)
			continue
		}

		tags, err := mp4.Read()
		if err != nil {
			mp4.Close()
			log.Printf("error reading tags from %s: %v", path, err)
			continue
		}

		// Inference: If track number is missing, try to parse it from the filename
		if tags.TrackNumber < 1 {
			d, t := inferTrackInfo(filepath.Base(path))
			if t > 0 {
				tags.TrackNumber = t
			}
			if d > 0 {
				tags.DiscNumber = d
			}
		}

		newTags, changed, err := c.UpdateFromMB(ctx, tags, releaseID)
		if err != nil {
			mp4.Close()
			log.Printf("error updating tags for %s: %v", path, err)
			continue
		}

		if changed {
			log.Printf("tagging %s (Track %d)...", f.Name(), newTags.TrackNumber)
			if !c.Config.DryRun {
				// Ensure Custom map is initialized
				if newTags.Custom == nil {
					newTags.Custom = make(map[string]string)
				}
				// Pass all custom tag keys to ensure they are saved
				customKeys := make([]string, 0, len(newTags.Custom))
				for k := range newTags.Custom {
					customKeys = append(customKeys, k)
				}
				if err := mp4.Write(newTags, customKeys); err != nil {
					log.Printf("error writing tags to %s: %v", path, err)
				}
			}
		}
		mp4.Close()

		if !c.Config.SkipMove {
			if err := c.rename(path, newTags); err != nil {
				log.Printf("error renaming %s: %v", path, err)
			}
		}
	}

	return nil
}

func inferTrackInfo(name string) (disc, track int16) {
	fields := strings.Fields(name)
	if len(fields) > 0 {
		firstUnit := fields[0]
		if strings.Contains(firstUnit, "-") {
			parts := strings.Split(firstUnit, "-")
			if len(parts) >= 2 {
				if d, err := strconv.Atoi(parts[0]); err == nil {
					disc = int16(d)
				}
				if t, err := strconv.Atoi(parts[1]); err == nil {
					track = int16(t)
				}
			}
		} else {
			if t, err := strconv.Atoi(firstUnit); err == nil {
				track = int16(t)
			}
		}
	}
	return
}
