package sa

import (
	"context"
	"fmt"
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
		out.Custom = make(map[string]string)
	} else {
		newCustom := make(map[string]string)
		for k, v := range out.Custom {
			newCustom[k] = v
		}
		out.Custom = newCustom
	}

	out.Album = release.Title
	out.AlbumSort = release.Title + " [discogs:" + strconv.Itoa(discogsID) + "]"
	out.Custom["Discogs Release Id"] = strconv.Itoa(discogsID)

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
		out.Custom["LABEL"] = strings.Join(labels, "; ")
		out.Custom["CATALOGNUMBER"] = strings.Join(cats, "; ")
	}

	if len(release.Genres) > 0 {
		out.Custom["Genre"] = strings.Join(release.Genres, "; ")
	}

	// Check if this Discogs release has a MusicBrainz ID in identifiers
	for _, id := range release.Identifiers {
		if id.Type == "barcode" && id.Value != "" {
			out.Custom["BARCODE"] = id.Value
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
