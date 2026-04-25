package sa

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/Sorrow446/go-mp4tag"
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

func (c *Curator) updateFromMB(in *mp4tag.MP4Tags, overrideID string) (*mp4tag.MP4Tags, bool, error) {
	if c.mb5query == nil {
		return in, false, fmt.Errorf("mb5 query not initialized")
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

	recordingID, _ := in.Custom["MusicBrainz Track Id"]

	if in.TrackNumber < 1 {
		log.Println("no track number, skipping")
		return in, false, nil
	}
	if in.DiscNumber < 1 {
		log.Println("no disc number, skipping")
		return in, false, nil
	}

	c.Stats.mu.Lock()
	if v, ok := c.Stats.BadQueries[releaseid]; ok && v {
		c.Stats.mu.Unlock()
		return in, false, fmt.Errorf("%s failed once already, skipping", releaseid)
	}
	c.Stats.mu.Unlock()
c.rl.Take()

var metadata mb5_metadata
var retryCount int
for {
	// Prepare query parameters for explicit "inc" request
	var params [1]*byte
	p1 := []byte("inc")
	params[0] = &p1[0]

	var values [1]*byte
	v1 := []byte("artists labels recordings release-groups url-rels artist-credits work-rels artist-rels work-level-rels")
	values[0] = &v1[0]

	// Query libmusicbrainz5
	metadata = mb5_query_query(c.mb5query, "release", releaseid, "", 1, unsafe.Pointer(&params), unsafe.Pointer(&values))
	if metadata != nil {
		break
	}

	lastCode := mb5_query_get_lasthttpcode(c.mb5query)
	if lastCode == 503 && retryCount < 3 {
		retryCount++
		log.Printf("MusicBrainz returned 503, retrying in 5 seconds (attempt %d/3)...", retryCount)
		time.Sleep(5 * time.Second)
		continue
	}

	c.Stats.mu.Lock()
	c.Stats.BadQueries[releaseid] = true
	c.Stats.mu.Unlock()

	var errbuf [256]byte
	mb5_query_get_lasterrormessage(c.mb5query, &errbuf[0], 256)
	cErr := strings.Trim(string(errbuf[:]), "\x00")
	return in, false, fmt.Errorf("query to MusicBrainz failed for %s (HTTP %d): %s", releaseid, lastCode, cErr)
}
defer mb5_metadata_delete(metadata)

	release := mb5_metadata_get_release(metadata)
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

	out.Album = mb5String(mb5_release_get_title, unsafe.Pointer(release))
	out.AlbumSort = out.Album + " [" + releaseid + "]"
	
	if disambig := mb5String(mb5_release_get_disambiguation, unsafe.Pointer(release)); disambig != "" {
		out.Comment = disambig
	}

	ac := mb5_release_get_artistcredit(release)
	if aa := c.fmtArtistCreditMB5(ac); aa != "" {
		out.AlbumArtist = aa
	}
	if asa := c.fmtArtistCreditSortMB5(ac); asa != "" && asa != out.AlbumArtist {
		out.AlbumArtistSort = asa
	}

	mediumList := mb5_release_get_mediumlist(release)
	out.DiscTotal = int16(mb5_medium_list_size(mediumList))

	var foundTrack mb5_track
	var foundMedium mb5_medium
	for i := 0; i < int(out.DiscTotal); i++ {
		m := mb5_medium_list_item(mediumList, i)
		if mb5_medium_get_position(m) == int(in.DiscNumber) {
			foundMedium = m
			tl := mb5_medium_get_tracklist(m)
			out.TrackTotal = int16(mb5_track_list_get_count(tl))
			for j := 0; j < int(out.TrackTotal); j++ {
				t := mb5_track_list_item(tl, j)
				if mb5_track_get_position(t) == int(in.TrackNumber) {
					foundTrack = t
					break
				}
			}
			break
		}
	}

	if foundTrack != nil {
		recording := mb5_track_get_recording(foundTrack)
		if recording != nil {
			if recordingID == "" {
				var buf [37]byte
				mb5_recording_get_id(unsafe.Pointer(recording), (*byte)(unsafe.Pointer(&buf[0])), 37)
				recordingID = strings.Trim(string(buf[:]), "\x00")
			}
			tac := mb5_recording_get_artistcredit(recording)
			if tac == nil {
				tac = mb5_track_get_artistcredit(foundTrack)
			}
			if artist := c.fmtArtistCreditMB5(tac); artist != "" {
				out.Artist = artist
			}
			
			title := mb5String(mb5_recording_get_title, unsafe.Pointer(recording))
			if title == "" {
				title = mb5String(mb5_track_get_title, unsafe.Pointer(foundTrack))
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

	rg := mb5_release_get_releasegroup(release)
	if rg != nil {
		origDate := mb5String(mb5_releasegroup_get_firstreleasedate, unsafe.Pointer(rg))
		if len(origDate) >= 4 {
			out.Custom["ORIGINALDATE"] = origDate
			out.Custom["ORIGINALYEAR"] = origDate[:4]
			fmt.Sscanf(origDate[:4], "%d", &out.Year)
		}
		if rgid := mb5String(mb5_releasegroup_get_id, unsafe.Pointer(rg)); rgid != "" {
			out.Custom["MusicBrainz Release Group Id"] = rgid
		}
		// Album Type is commented out in mb5.go
	}

	if date := mb5String(mb5_release_get_date, unsafe.Pointer(release)); date != "" {
		out.Date = date
	}

	if asin := mb5String(mb5_release_get_asin, unsafe.Pointer(release)); asin != "" {
		out.Custom["ASIN"] = asin
	}
	if barcode := mb5String(mb5_release_get_barcode, unsafe.Pointer(release)); barcode != "" {
		out.Custom["BARCODE"] = barcode
	}
	
	liList := mb5_release_get_labelinfolist(release)
	if liList != nil {
		if label := c.fmtLabelsMB5(liList); label != "" {
			out.Custom["LABEL"] = label
		}
		if cat := c.fmtCatalogNumbersMB5(liList); cat != "" {
			out.Custom["CATALOGNUMBER"] = cat
		}
	}

	countryCode := mb5String(mb5_release_get_country, unsafe.Pointer(release))
	if countryCode != "" {
		if countryName := resolveCountry(countryCode); countryName != "" {
			out.Custom["Country"] = countryName
		}
		out.Custom["MusicBrainz Album Release Country"] = countryCode
	}
	
	if foundMedium != nil {
		if format := mb5String(mb5_medium_get_format, unsafe.Pointer(foundMedium)); format != "" {
			out.Custom["MEDIA"] = mediumFormat(format)
		}
	}

	// LANGUAGE and SCRIPT are commented out in mb5.go
	if aids := c.joinArtistIDsMB5(ac); aids != "" {
		out.Custom["MusicBrainz Album Artist Id"] = aids
	}
	out.Custom["MusicBrainz Album Id"] = releaseid
	out.Custom["MusicBrainz Track Id"] = recordingID

	if c.Config.Debug {
		temp := out
		temp.Pictures = nil
		log.Printf("AFTER MB UPDATE: %# v", pretty.Formatter(temp))
	}

	return &out, true, nil
}

func (c *Curator) processRelations(recording mb5_recording, custom map[string]string, out *mp4tag.MP4Tags) {
	rll := mb5_recording_get_relationlistlist(recording)
	if rll == nil {
		return
	}

	roles := map[string]string{
		"composer":      "COMPOSER",
		"lyricist":      "LYRICIST",
		"producer":      "PRODUCER",
		"engineer":      "ENGINEER",
		"mixer":         "MIXER",
		"remixer":       "REMIXER",
		"writer":        "WRITER",
		"arranger":      "ARRANGER",
		"conductor":     "CONDUCTOR",
	}

	credits := make(map[string][]string)

	rllSize := mb5_relationlist_list_size(rll)
	for i := 0; i < rllSize; i++ {
		rl := mb5_relationlist_list_item(rll, i)
		rlSize := mb5_relation_list_size(rl)
		for j := 0; j < rlSize; j++ {
			rel := mb5_relation_list_item(rl, j)
			relType := mb5String(mb5_relation_get_type, unsafe.Pointer(rel))
			
			if role, ok := roles[relType]; ok {
				artist := mb5_relation_get_artist(rel)
				if artist != nil {
					name := mb5String(mb5_artist_get_name, unsafe.Pointer(artist))
					if name != "" {
						credits[role] = append(credits[role], name)
					}
				}
			}

			if relType == "performance" {
				work := mb5_relation_get_work(rel)
				if work != nil {
					if wid := mb5String(mb5_work_get_id, unsafe.Pointer(work)); wid != "" {
						custom["MusicBrainz Work Id"] = wid
					}
					if wtitle := mb5String(mb5_work_get_title, unsafe.Pointer(work)); wtitle != "" {
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

func (c *Curator) fmtISRCsMB5(recording mb5_recording) string {
	isrcList := mb5_recording_get_isrclist(recording)
	if isrcList == nil {
		return ""
	}
	size := mb5_isrc_list_size(isrcList)
	var isrcs []string
	for i := 0; i < size; i++ {
		isrc := mb5_isrc_list_item(isrcList, i)
		id := mb5String(mb5_isrc_get_id, unsafe.Pointer(isrc))
		if id != "" {
			isrcs = append(isrcs, id)
		}
	}
	return strings.Join(isrcs, ", ")
}

func (c *Curator) fmtArtistCreditMB5(ac mb5_artist_credit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5_artistcredit_get_namecreditlist(ac)
	count := mb5_namecredit_list_get_count(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5_namecredit_list_item(ncl, i)
		name := mb5String(mb5_namecredit_get_name, unsafe.Pointer(nc))
		if name == "" {
			artist := mb5_namecredit_get_artist(nc)
			name = mb5String(mb5_artist_get_name, unsafe.Pointer(artist))
		}
		join := mb5String(mb5_namecredit_get_joinphrase, unsafe.Pointer(nc))
		s += name + join
	}
	return s
}

func (c *Curator) fmtArtistCreditSortMB5(ac mb5_artist_credit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5_artistcredit_get_namecreditlist(ac)
	count := mb5_namecredit_list_get_count(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5_namecredit_list_item(ncl, i)
		artist := mb5_namecredit_get_artist(nc)
		name := mb5String(mb5_artist_get_sortname, unsafe.Pointer(artist))
		join := mb5String(mb5_namecredit_get_joinphrase, unsafe.Pointer(nc))
		s += name + join
	}
	return s
}

func (c *Curator) fmtArtistListMB5(ac mb5_artist_credit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5_artistcredit_get_namecreditlist(ac)
	count := mb5_namecredit_list_get_count(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5_namecredit_list_item(ncl, i)
		artist := mb5_namecredit_get_artist(nc)
		name := mb5String(mb5_artist_get_name, unsafe.Pointer(artist))
		if i > 0 {
			s += ", "
		}
		s += name
	}
	return s
}

func (c *Curator) joinArtistIDsMB5(ac mb5_artist_credit) string {
	if ac == nil {
		return ""
	}
	ncl := mb5_artistcredit_get_namecreditlist(ac)
	count := mb5_namecredit_list_get_count(ncl)
	var s string
	for i := 0; i < count; i++ {
		nc := mb5_namecredit_list_item(ncl, i)
		artist := mb5_namecredit_get_artist(nc)
		id := mb5String(mb5_artist_get_id, unsafe.Pointer(artist))
		if i > 0 {
			s += ","
		}
		s += id
	}
	return s
}

func (c *Curator) fmtLabelsMB5(liList mb5_label_info_list) string {
	count := mb5_labelinfo_list_size(liList)
	var labels []string
	seen := make(map[string]bool)
	for i := 0; i < count; i++ {
		li := mb5_labelinfo_list_item(liList, i)
		label := mb5_labelinfo_get_label(li)
		if label != nil {
			name := mb5String(mb5_label_get_name, unsafe.Pointer(label))
			if name != "" && !seen[name] {
				labels = append(labels, name)
				seen[name] = true
			}
		}
	}
	return strings.Join(labels, "; ")
}

func (c *Curator) fmtCatalogNumbersMB5(liList mb5_label_info_list) string {
	count := mb5_labelinfo_list_size(liList)
	var cats []string
	seen := make(map[string]bool)
	for i := 0; i < count; i++ {
		li := mb5_labelinfo_list_item(liList, i)
		cat := mb5String(mb5_labelinfo_get_catalognumber, unsafe.Pointer(li))
		if cat != "" && !seen[cat] {
			cats = append(cats, cat)
			seen[cat] = true
		}
	}
	return strings.Join(cats, "; ")
}

func (c *Curator) updateFromDiscID(in *mp4tag.MP4Tags, discid string) (*mp4tag.MP4Tags, bool, error) {
	if c.mb5query == nil {
		return in, false, fmt.Errorf("mb5 query not initialized")
	}

	c.rl.Take()

	var relList mb5_release_list
	var retryCount int
	for {
		relList = mb5_query_lookup_discid(c.mb5query, mb5_discid(discid))
		if relList != nil {
			break
		}

		lastCode := mb5_query_get_lasthttpcode(c.mb5query)
		if lastCode == 503 && retryCount < 3 {
			retryCount++
			log.Printf("MusicBrainz returned 503 (discid), retrying in 5 seconds (attempt %d/3)...", retryCount)
			time.Sleep(5 * time.Second)
			continue
		}

		return in, false, fmt.Errorf("discid lookup failed for %s (HTTP %d)", discid, lastCode)
	}
	defer mb5_release_list_delete(relList)

	if mb5_release_list_size(relList) == 0 {
		return in, false, fmt.Errorf("no releases found for discid %s", discid)
	}

	// For now, we take the first release and use its ID
	rel := mb5_release_list_item(relList, 0)
	var relBuf [256]byte
	mb5_release_get_id(unsafe.Pointer(rel), &relBuf[0], 256)
	releaseID := string(relBuf[:cStringLen(relBuf[:])])

	if c.Config.Debug {
		log.Printf("resolved discid %s to release %s\n", discid, releaseID)
	}

	// Now call the standard update with this release ID
	return c.updateFromMB(in, releaseID)
}

func (c *Curator) updateFromTOC(in *mp4tag.MP4Tags, toc string) (*mp4tag.MP4Tags, bool, error) {
	if c.mb5query == nil {
		return in, false, fmt.Errorf("mb5 query not initialized")
	}

	c.rl.Take()

	var relList mb5_release_list
	var retryCount int
	for {
		relList = mb5_query_lookup_toc(c.mb5query, toc)
		if relList != nil {
			break
		}

		lastCode := mb5_query_get_lasthttpcode(c.mb5query)
		if lastCode == 503 && retryCount < 3 {
			retryCount++
			log.Printf("MusicBrainz returned 503 (toc), retrying in 5 seconds (attempt %d/3)...", retryCount)
			time.Sleep(5 * time.Second)
			continue
		}

		return in, false, fmt.Errorf("toc lookup failed for %s (HTTP %d)", toc, lastCode)
	}
	defer mb5_release_list_delete(relList)

	if mb5_release_list_size(relList) == 0 {
		return in, false, fmt.Errorf("no releases found for toc %s", toc)
	}

	// For now, we take the first release and use its ID
	rel := mb5_release_list_item(relList, 0)
	var relBuf [256]byte
	mb5_release_get_id(unsafe.Pointer(rel), &relBuf[0], 256)
	releaseID := string(relBuf[:cStringLen(relBuf[:])])

	if c.Config.Debug {
		log.Printf("resolved toc %s to release %s\n", toc, releaseID)
	}

	// Now call the standard update with this release ID
	return c.updateFromMB(in, releaseID)
}

func (c *Curator) acoustIDLookup(fingerprint string, duration int) (string, error) {
	if c.Config.AcoustIDKey == "" {
		return "", fmt.Errorf("AcoustID key not set")
	}

	url := fmt.Sprintf("https://api.acoustid.org/v2/lookup?client=%s&meta=releasegroups+releases&duration=%d&fingerprint=%s",
		c.Config.AcoustIDKey, duration, fingerprint)

	resp, err := http.Get(url)
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

	return "", fmt.Errorf("no high-confidence matches found in AcoustID")
}
