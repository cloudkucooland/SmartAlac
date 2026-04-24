package sa

import (
	"fmt"
	"log"
	"strings"
	"unsafe"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/kr/pretty"
)

func (c *Curator) updateFromMB(in *mp4tag.MP4Tags) (*mp4tag.MP4Tags, bool, error) {
	if c.mb5query == nil {
		return in, false, fmt.Errorf("mb5 query not initialized")
	}

	if c.Config.Debug {
		log.Printf("%# v", pretty.Formatter(in))
	}

	releaseid, ok := in.Custom["MusicBrainz Album Id"]
	if !ok {
		log.Println("no release ID, skipping")
		return in, false, nil
	}

	recordingID, ok := in.Custom["MusicBrainz Track Id"]
	if !ok {
		log.Println("no recordingID, skipping")
		return in, false, nil
	}

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

	// Prepare query parameters for explicit "inc" request
	// Expanded to include relations and works
	var params [1]*byte
	p1 := []byte("inc")
	params[0] = &p1[0]

	var values [1]*byte
	v1 := []byte("artists labels recordings release-groups url-rels artist-credits work-rels artist-rels work-level-rels")
	values[0] = &v1[0]

	// Query libmusicbrainz5
	metadata := mb5_query_query(c.mb5query, "release", releaseid, "", 1, unsafe.Pointer(&params), unsafe.Pointer(&values))
	if metadata == nil {
		c.Stats.mu.Lock()
		c.Stats.BadQueries[releaseid] = true
		c.Stats.mu.Unlock()

		var errbuf [256]byte
		mb5_query_get_lasterrormessage(c.mb5query, &errbuf[0], 256)
		cErr := strings.Trim(string(errbuf[:]), "\x00")
		return in, false, fmt.Errorf("query to MusicBrainz failed for %s: %s", releaseid, cErr)
	}
	defer mb5_metadata_delete(metadata)

	release := mb5_metadata_get_release(metadata)
	if release == nil {
		return in, false, fmt.Errorf("no release in metadata for %s", releaseid)
	}

	out := mp4tag.MP4Tags{
		ItunesAdvisory: 0,
		ItunesAlbumID:  -1,
		ItunesArtistID: -1,
	}
	out.Custom = make(map[string]string)

	// Copy preserved customs
	copyCustoms := []string{"KEY", "MOOD", "URL_LYRICS_SITE", "VINYLDIGITIZER", "URL_DISCOGS_ARTIST_SITE", "DIGITIZE_DATE", "DIGITIZE_INFO", "MusicBrainz Disc Id", "initialkey", "MusicBrainz Release Track Id"}
	for _, v := range copyCustoms {
		if val, ok := in.Custom[v]; ok {
			out.Custom[v] = val
		}
	}

	out.Album = mb5String(mb5_release_get_title, unsafe.Pointer(release))
	out.AlbumSort = out.Album + " [" + releaseid + "]"
	
	if disambig := mb5String(mb5_release_get_disambiguation, unsafe.Pointer(release)); disambig != "" {
		out.Comment = disambig
	} else {
		out.Comment = in.Comment
	}

	ac := mb5_release_get_artistcredit(release)
	aa := c.fmtArtistCreditMB5(ac)
	asa := c.fmtArtistCreditSortMB5(ac)
	out.AlbumArtist = aa
	if aa != asa {
		out.AlbumArtistSort = asa
	}

	out.BPM = in.BPM
	out.DiscNumber = in.DiscNumber
	
	mediumList := mb5_release_get_mediumlist(release)
	out.DiscTotal = int16(mb5_medium_list_size(mediumList))

	// Find track and recording
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
			tac := mb5_recording_get_artistcredit(recording)
			if tac == nil {
				tac = mb5_track_get_artistcredit(foundTrack)
			}
			out.Artist = c.fmtArtistCreditMB5(tac)
			
			title := mb5String(mb5_recording_get_title, unsafe.Pointer(recording))
			if title == "" {
				title = mb5String(mb5_track_get_title, unsafe.Pointer(foundTrack))
			}
			out.Title = title

			out.Custom["ARTISTS"] = c.fmtArtistListMB5(tac)
			out.Custom["MusicBrainz Artist Id"] = c.joinArtistIDsMB5(tac)

			// Advanced: Relations and ISRCs
			c.processRelations(recording, out.Custom, &out)
			out.Custom["ISRC"] = c.fmtISRCsMB5(recording)
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
		out.Custom["MusicBrainz Release Group Id"] = mb5String(mb5_releasegroup_get_id, unsafe.Pointer(rg))
		out.Custom["MusicBrainz Album Type"] = strings.ToLower(mb5String(mb5_releasegroup_get_type, unsafe.Pointer(rg)))
	}

	out.CustomGenre = in.CustomGenre
	out.Date = mb5String(mb5_release_get_date, unsafe.Pointer(release))
	out.Lyrics = in.Lyrics
	out.TrackNumber = in.TrackNumber

	if asin := mb5String(mb5_release_get_asin, unsafe.Pointer(release)); asin != "" {
		out.Custom["ASIN"] = asin
	}
	out.Custom["BARCODE"] = mb5String(mb5_release_get_barcode, unsafe.Pointer(release))
	
	liList := mb5_release_get_labelinfolist(release)
	if liList != nil {
		out.Custom["LABEL"] = c.fmtLabelsMB5(liList)
		out.Custom["CATALOGNUMBER"] = c.fmtCatalogNumbersMB5(liList)
	}

	countryCode := mb5String(mb5_release_get_country, unsafe.Pointer(release))
	out.Custom["Country"] = resolveCountry(countryCode)
	out.Custom["MusicBrainz Album Release Country"] = countryCode
	out.Custom["MusicBrainz Album Status"] = strings.ToLower(mb5String(mb5_release_get_status, unsafe.Pointer(release)))
	
	if foundMedium != nil {
		out.Custom["MEDIA"] = mb5String(mb5_medium_get_format, unsafe.Pointer(foundMedium))
	}

	out.Custom["LANGUAGE"] = mb5String(mb5_release_get_language, unsafe.Pointer(release))
	out.Custom["SCRIPT"] = mb5String(mb5_release_get_script, unsafe.Pointer(release))
	out.Custom["MusicBrainz Album Artist Id"] = c.joinArtistIDsMB5(ac)
	out.Custom["MusicBrainz Album Id"] = releaseid
	out.Custom["MusicBrainz Track Id"] = recordingID

	// Copy remaining unhandled implements if they still exist in input
	copyCustomsNotImpl := []string{"MusicBrainz Work Id", "WORK", "LYRICIST", "PRODUCER", "ENGINEER", "MIXER", "REMIXER", "WRITER", "ARRANGER", "DISCSUBTITLE"}
	for _, v := range copyCustomsNotImpl {
		if val, ok := in.Custom[v]; ok {
			if _, exists := out.Custom[v]; !exists {
				out.Custom[v] = val
			}
		}
	}
	out.Composer = in.Composer // We might overwrite this in processRelations

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
			
			// Handle Artist relations
			if role, ok := roles[relType]; ok {
				artist := mb5_relation_get_artist(rel)
				if artist != nil {
					name := mb5String(mb5_artist_get_name, unsafe.Pointer(artist))
					credits[role] = append(credits[role], name)
				}
			}

			// Handle Work relations
			if relType == "performance" {
				work := mb5_relation_get_work(rel)
				if work != nil {
					custom["MusicBrainz Work Id"] = mb5String(mb5_work_get_id, unsafe.Pointer(work))
					custom["WORK"] = mb5String(mb5_work_get_title, unsafe.Pointer(work))
				}
			}
		}
	}

	for role, names := range credits {
		val := strings.Join(names, ", ")
		custom[role] = val
		if role == "COMPOSER" {
			out.Composer = val
		}
		if role == "CONDUCTOR" {
			out.Conductor = val
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
