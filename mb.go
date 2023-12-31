package sa

import (
	"fmt"
	"log"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/kr/pretty"
	"github.com/michiwend/gomusicbrainz"
	// if the caches are too big, try one of these (or others)
	// git.mills.io/prologic/bitcask
	// github.com/dgraph-io/ristretto
	// github.com/peterbourgon/diskv
)

var client *gomusicbrainz.WS2Client
var releases map[string]*gomusicbrainz.Release

func init() {
	client, _ = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"SmartAlac",
		"0.0.0",
		"http://github.com/cloudkucooland/SmartAlac")

	releases = make(map[string]*gomusicbrainz.Release)
}

func updateFromMB(in *mp4tag.MP4Tags) (*mp4tag.MP4Tags, bool, error) {
	out := mp4tag.MP4Tags{
		ItunesAdvisory: 0,
		ItunesAlbumID:  -1,
		ItunesArtistID: -1,
	}

	out.Custom = make(map[string]string)
	copyCustoms := []string{"KEY", "MOOD", "URL_LYRICS_SITE", "VINYLDIGITIZER", "URL_DISCOGS_ARTIST_SITE"}
	for _, v := range copyCustoms {
		if _, ok := in.Custom[v]; ok {
			out.Custom[v] = in.Custom[v]
		}
	}

	tid, ok := in.Custom["MusicBrainz Track Id"]
	if !ok {
		log.Println("no track ID, skipping")
		return in, false, nil
	}

	releaseid, ok := in.Custom["MusicBrainz Album Id"]
	if !ok {
		log.Println("no release ID, skipping")
		return in, false, nil
	}

	release, ok := releases[releaseid]
	if !ok {
		var err error
		release, err = client.LookupRelease(gomusicbrainz.MBID(releaseid), "artist-credits", "recordings", "release-groups", "media", "isrcs", "release-rels", "release-group-rels", "url-rels", "labels")
		if err != nil {
			log.Println(err.Error())
			return in, false, err
		}
		releases[releaseid] = release
		// log.Printf("%# v\n", pretty.Formatter(release))
	}

	// per release items
	out.Album = release.Title
	out.AlbumSort = release.Title
	out.AlbumArtist = fmtArtistCredit(release.ArtistCredit.NameCredits)
	out.AlbumArtistSort = fmtArtistCreditSort(release.ArtistCredit.NameCredits)
	out.BPM = in.BPM
	out.DiscNumber = in.DiscNumber
	out.DiscTotal = getMediumCount(release)
	// out.Sort:
	out.TrackTotal = getMediumTrackCount(release, in.DiscNumber)
	// out.Publisher =
	out.Year = int32(release.ReleaseGroup.FirstReleaseDate.Year())

	// per track items
	medium := getMedium(release, in.DiscNumber)
	track := getTrack(release, gomusicbrainz.MBID(tid))
	out.Artist = fmtArtistCredit(track.Recording.ArtistCredit.NameCredits)
	out.Comment = in.Comment
	// out.Composer:
	// out.ComposerSort:
	// out.Conductor:
	// out.Copyright:

	out.CustomGenre = in.CustomGenre
	// out.Date:
	// out.Description:
	// out.Director:
	// out.Genre:0
	out.Lyrics = in.Lyrics
	// out.Narrator:
	out.Title = track.Recording.Title
	// out.TitleSort:
	out.TrackNumber = in.TrackNumber

	out.Custom["ARTISTS"] = fmtArtistList(release.ArtistCredit.NameCredits)
	out.Custom["BARCODE"] = release.Barcode
	out.Custom["CATALOGNUMBER"] = fmtCatalogNumber(release.LabelInfos)
	out.Custom["Country"] = release.CountryCode // need to convert to country name
	out.Custom["LABEL"] = fmtLabel(release.LabelInfos)
	out.Custom["LANGUAGE"] = release.TextRepresentation.Language
	out.Custom["MEDIA"] = medium.Format
	out.Custom["MusicBrainz Album Artist Id"] = joinArtistIDs(release.ArtistCredit.NameCredits)
	out.Custom["MusicBrainz Album Id"] = releaseid
	out.Custom["MusicBrainz Album Release Country"] = release.CountryCode
	out.Custom["MusicBrainz Album Type"] = strings.ToLower(release.ReleaseGroup.Type)
	out.Custom["MusicBrainz Artist Id"] = joinArtistIDs(track.Recording.ArtistCredit.NameCredits)
	// out.Custom["MusicBrainz Disc Id"] :YeZdB4dY7lQIC0ZSgkbe2R4HXAs-
	out.Custom["MusicBrainz Release Group Id"] = string(release.ReleaseGroup.ID)
	out.Custom["MusicBrainz Release Track Id"] = tid
	out.Custom["MusicBrainz Track Id"] = tid
	out.Custom["ORIGINALDATE"] = fmt.Sprintf("%d", release.ReleaseGroup.FirstReleaseDate.Year())
	out.Custom["ORIGINALYEAR"] = fmt.Sprintf("%d", release.ReleaseGroup.FirstReleaseDate.Year())
	out.Custom["RELEASESTATUS"] = strings.ToLower(release.Status)
	out.Custom["SCRIPT"] = release.TextRepresentation.Script

	// out.Custom["URL_DISCOGS_RELEASE_SITE"] = release.Relations.Url.RelationAbstract.Target

	log.Printf("%# v\n", pretty.Formatter(pretty.Diff(in, &out)))

	return &out, true, nil
}

func getMediumTrackCount(r *gomusicbrainz.Release, discnumber int16) int16 {
	for _, m := range r.Mediums {
		if int16(m.Position) == discnumber {
			return int16(len(m.Tracks))
		}
	}
	return 0
}

func getMediumCount(r *gomusicbrainz.Release) int16 {
	return int16(len(r.Mediums))
}

func getMedium(r *gomusicbrainz.Release, discnumber int16) *gomusicbrainz.Medium {
	for _, m := range r.Mediums {
		if int16(m.Position) == discnumber {
			return m
		}
	}
	return nil
}

func fmtArtistCredit(a []gomusicbrainz.NameCredit) string {
	var s string
	for k, v := range a {
		if k > 0 {
			s += " "
		}
		s += v.Artist.Name
	}
	// log.Printf("fmtArtistCredit: %s", s)
	return s
}

func fmtArtistList(a []gomusicbrainz.NameCredit) string {
	var s string
	for k, v := range a {
		if k > 0 {
			s += ", "
		}
		s += v.Artist.Name
	}
	// log.Printf("fmtArtistList: %s", s)
	return s
}

func fmtArtistCreditSort(a []gomusicbrainz.NameCredit) string {
	var s string
	for k, v := range a {
		if k > 0 {
			s += ", "
		}
		s += v.Artist.SortName
	}
	// log.Printf("fmtArtistCreditSort: %s", s)
	return s
}

func joinArtistIDs(a []gomusicbrainz.NameCredit) string {
	var s string
	for k, v := range a {
		if k > 0 {
			s += ","
		}
		s += string(v.Artist.ID)
	}
	return s
}

func getTrack(r *gomusicbrainz.Release, trackID gomusicbrainz.MBID) *gomusicbrainz.Track {
	for _, m := range r.Mediums {
		for _, t := range m.Tracks {
			if t.ID == trackID {
				return t
			}
		}
	}
	log.Println("unable to find matching trackID")
	return nil
}

func fmtCatalogNumber(l []gomusicbrainz.LabelInfo) string {
	var s string
	for i, li := range l {
		if i > 0 {
			s += ", "
		}
		s += li.CatalogNumber
	}
	return s
}

func fmtLabel(l []gomusicbrainz.LabelInfo) string {
	var s string
	for i, li := range l {
		if i > 0 {
			s += ", "
		}
		s += li.Label.Name
	}
	return s
}
