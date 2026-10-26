package sa

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/irlndts/go-discogs"
)

func (c *Curator) UpdateFromDiscogs(ctx context.Context, in *mp4tag.MP4Tags, discogsID int) (*mp4tag.MP4Tags, bool, error) {
	if c.dgClient == nil {
		return in, false, fmt.Errorf("discogs client not initialized")
	}

	client := c.dgClient.(discogs.Discogs)

	select {
	case <-ctx.Done():
		return in, false, ctx.Err()
	default:
	}

	c.dgRL.Take()
	release, err := client.Release(discogsID)
	if err != nil {
		return in, false, fmt.Errorf("discogs lookup failed for %d: %w", discogsID, err)
	}

	out := *in
	if out.Custom == nil {
		out.Custom = make(map[string][]string)
	} else {
		newCustom := make(map[string][]string)
		for k, v := range out.Custom {
			vCopy := make([]string, len(v))
			copy(vCopy, v)
			newCustom[k] = vCopy
		}
		out.Custom = newCustom
	}

	out.Album = release.Title
	out.AlbumSort = release.Title + " [discogs:" + strconv.Itoa(discogsID) + "]"
	out.Custom["Discogs Release Id"] = []string{strconv.Itoa(discogsID)}

	if len(release.Artists) > 0 {
		var artists []string
		for _, a := range release.Artists {
			artists = append(artists, a.Name)
		}
		out.AlbumArtist = strings.Join(artists, ", ")
	}

	// Match track
	var foundTrack *discogs.Track
	for _, t := range release.Tracklist {
		if t.Position == strconv.Itoa(int(in.TrackNumber)) {
			foundTrack = &t
			break
		}
		// Also check if position is just the number
		if t.Position == fmt.Sprintf("%02d", in.TrackNumber) {
			foundTrack = &t
			break
		}
	}

	if foundTrack != nil {
		out.Title = foundTrack.Title
		if len(foundTrack.Artists) > 0 {
			var artists []string
			for _, a := range foundTrack.Artists {
				artists = append(artists, a.Name)
			}
			out.Artist = strings.Join(artists, ", ")
		} else {
			out.Artist = out.AlbumArtist
		}
	}

	if release.Year != 0 {
		// out.Year is int32 in go-mp4tag, release.Year is int
		out.Year = int32(release.Year)
		out.Date = strconv.Itoa(release.Year)
	}

	if len(release.Labels) > 0 {
		var labels []string
		var cats []string
		for _, l := range release.Labels {
			labels = append(labels, l.Name)
			if l.Catno != "" {
				cats = append(cats, l.Catno)
			}
		}
		out.Custom["LABEL"] = []string{strings.Join(labels, "; ")}
		out.Custom["CATALOGNUMBER"] = []string{strings.Join(cats, "; ")}
	}

	if len(release.Genres) > 0 {
		out.Custom["Genre"] = []string{strings.Join(release.Genres, "; ")}
	}

	// Check if this Discogs release has a MusicBrainz ID in identifiers
	for _, id := range release.Identifiers {
		if id.Type == "barcode" && id.Value != "" {
			out.Custom["BARCODE"] = []string{id.Value}
		}
	}

	changed := !tagsEquivalent(in, &out)
	return &out, changed, nil
}
func (c *Curator) SearchDiscogs(ctx context.Context, artist, album, barcode string) ([]discogs.Result, error) {
	if c.dgClient == nil {
		return nil, fmt.Errorf("discogs client not initialized")
	}

	client := c.dgClient.(discogs.Discogs)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c.dgRL.Take()
	request := discogs.SearchRequest{
		Artist:       artist,
		ReleaseTitle: album,
		Barcode:      barcode,
		Type:         "release",
	}

	res, err := client.Search(request)
	if err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (c *Curator) SearchDiscogsByTrackDurations(ctx context.Context, trackDurations []int) ([]discogs.Result, error) {
	if c.dgClient == nil {
		return nil, fmt.Errorf("discogs client not initialized")
	}

	trackCount := len(trackDurations)
	if trackCount == 0 {
		return nil, fmt.Errorf("no track durations provided")
	}

	client := c.dgClient.(discogs.Discogs)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Step 1: Search Discogs by Format (CD) and total Track Count query
	c.dgRL.Take()
	request := discogs.SearchRequest{
		Format: "CD",
		Type:   "release",
		Q:      fmt.Sprintf("%d tracks", trackCount),
	}

	res, err := client.Search(request)
	if err != nil {
		return nil, fmt.Errorf("discogs track count search failed: %w", err)
	}

	if len(res.Results) == 0 {
		return nil, nil
	}

	// Step 2: Fetch individual releases to inspect tracklist durations
	var matchedResults []discogs.Result
	for _, result := range res.Results {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		c.dgRL.Take()
		release, err := client.Release(result.ID)
		if err != nil {
			continue // Skip unresolvable release entries
		}

		// Filter against the target duration profile
		if matchDiscogsTracklist(release.Tracklist, trackDurations) {
			matchedResults = append(matchedResults, result)
		}
	}

	return matchedResults, nil
}

// matchDiscogsTracklist parses MM:SS / HH:MM:SS duration strings in Discogs
// tracklists and checks if they match the expected track durations (+/- 3s tolerance per track).
func matchDiscogsTracklist(tracklist []discogs.Track, targetDurations []int) bool {
	// Filter out non-audio tracks (headers, indexes)
	var audioTracks []discogs.Track
	for _, t := range tracklist {
		if t.Duration != "" || t.Position != "" {
			audioTracks = append(audioTracks, t)
		}
	}

	if len(audioTracks) != len(targetDurations) {
		return false
	}

	matchedCount := 0
	for i, t := range audioTracks {
		discogsSecs := parseDiscogsDuration(t.Duration)
		if discogsSecs == 0 {
			// Track lacks duration metadata on Discogs side; give benefit of the doubt
			continue
		}

		// Allow a 3-second tolerance per track to handle frame rounding/pregap differences
		if math.Abs(float64(discogsSecs-targetDurations[i])) <= 3.0 {
			matchedCount++
		}
	}

	// Require at least 80% of tracks with valid durations to match
	return float64(matchedCount)/float64(len(targetDurations)) >= 0.8
}

// parseDiscogsDuration converts "MM:SS" or "HH:MM:SS" strings into total seconds.
func parseDiscogsDuration(dur string) int {
	parts := strings.Split(dur, ":")
	if len(parts) == 2 {
		min, _ := strconv.Atoi(parts[0])
		sec, _ := strconv.Atoi(parts[1])
		return min*60 + sec
	} else if len(parts) == 3 {
		hr, _ := strconv.Atoi(parts[0])
		min, _ := strconv.Atoi(parts[1])
		sec, _ := strconv.Atoi(parts[2])
		return hr*3600 + min*60 + sec
	}
	return 0
}
