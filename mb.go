package sa

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/ratelimit"

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
var rl ratelimit.Limiter

func init() {
	client, _ = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"SmartAlac",
		"0.0.0",
		"http://github.com/cloudkucooland/SmartAlac")

	releases = make(map[string]*gomusicbrainz.Release)

	rl = ratelimit.New(1)
}

func updateFromMB(in *mp4tag.MP4Tags) (*mp4tag.MP4Tags, bool, error) {
	out := mp4tag.MP4Tags{
		ItunesAdvisory: 0,
		ItunesAlbumID:  -1,
		ItunesArtistID: -1,
	}

	out.Custom = make(map[string]string)

	// Things generated elsewhere, just blindly copy them over for now
	copyCustoms := []string{"KEY", "MOOD", "URL_LYRICS_SITE", "VINYLDIGITIZER", "URL_DISCOGS_ARTIST_SITE", "DIGITIZE_DATE", "DIGITIZE_INFO", "MusicBrainz Disc Id"}
	for _, v := range copyCustoms {
		if _, ok := in.Custom[v]; ok {
			out.Custom[v] = in.Custom[v]
		}
	}

	releaseid, ok := in.Custom["MusicBrainz Album Id"]
	if !ok {
		log.Println("no release ID, skipping")
		return in, false, nil
	}
	log.Printf("release ID: %s", releaseid)

	tid, ok := in.Custom["MusicBrainz Track Id"]
	if !ok {
		log.Println("no track ID, skipping")
		return in, false, nil
	}
	if in.TrackNumber < 1 {
		log.Println("no track number, skipping")
		return in, false, nil
	}

	release, ok := releases[releaseid]
	if !ok {
		if v, ok := stats.badqueries[releaseid]; ok && v {
			err := fmt.Errorf("%s failed once already, skipping", releaseid)
			log.Println(err.Error())
			return in, false, err
		}

		rl.Take()
		var err error
		// release, err = client.LookupRelease(gomusicbrainz.MBID(releaseid), "artist-credits", "recordings", "release-groups", "media", "isrcs", "release-rels", "release-group-rels", "url-rels", "labels", "artists", "work-rels")
		release, err = client.LookupRelease(gomusicbrainz.MBID(releaseid), "artist-credits", "recordings", "release-groups", "media", "url-rels", "labels", "artists")
		if err != nil {
			stats.badqueries[releaseid] = true
			log.Printf("query to MusicBrainz failed for %s: %s\n", releaseid, err.Error())
			return in, false, err
		}

		if debug {
			log.Printf("%# v\n", pretty.Formatter(release))
		}

		releases[releaseid] = release
	}

	// per release items
	out.Album = release.Title
	out.AlbumSort = release.Title + " [" + releaseid + "]"

	// do not save AlbumArtistSort if not different than AlbumArtist
	aa := fmtArtistCredit(release.ArtistCredit.NameCredits)
	asa := fmtArtistCreditSort(release.ArtistCredit.NameCredits)
	out.AlbumArtist = aa
	if aa != asa {
		out.AlbumArtistSort = asa
	}
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
	if track == nil {
		return in, false, fmt.Errorf("no matching track in MB data")
	}
	out.Artist = fmtArtistCredit(track.Recording.ArtistCredit.NameCredits)
	out.Comment = in.Comment
	// out.Conductor:
	// out.Copyright:
	out.CustomGenre = in.CustomGenre
	out.Date = formatDate(release.Date)
	// out.Description:
	// out.Director:
	// out.Genre:
	out.Lyrics = in.Lyrics
	// out.Narrator:
	out.Title = track.Recording.Title
	// out.TitleSort:
	out.TrackNumber = in.TrackNumber

	out.Custom["ARTISTS"] = fmtArtistList(track.Recording.ArtistCredit.NameCredits)
	if release.Asin != "" {
		out.Custom["ASIN"] = release.Asin
	}
	out.Custom["BARCODE"] = release.Barcode
	out.Custom["CATALOGNUMBER"] = fmtCatalogNumber(release.LabelInfos)
	out.Custom["Country"] = resolveCountry(release.CountryCode)
	out.Custom["LABEL"] = fmtLabel(release.LabelInfos)
	out.Custom["LANGUAGE"] = release.TextRepresentation.Language
	out.Custom["MEDIA"] = mediumFormat(medium.Format)
	out.Custom["MusicBrainz Album Artist Id"] = joinArtistIDs(release.ArtistCredit.NameCredits)
	out.Custom["MusicBrainz Album Id"] = releaseid
	out.Custom["MusicBrainz Album Release Country"] = release.CountryCode
	out.Custom["MusicBrainz Album Type"] = strings.ToLower(release.ReleaseGroup.Type)
	out.Custom["MusicBrainz Artist Id"] = joinArtistIDs(track.Recording.ArtistCredit.NameCredits)
	out.Custom["MusicBrainz Release Group Id"] = string(release.ReleaseGroup.ID)
	out.Custom["MusicBrainz Release Track Id"] = tid
	out.Custom["MusicBrainz Track Id"] = tid
	out.Custom["ORIGINALDATE"] = formatDate(release.ReleaseGroup.FirstReleaseDate)
	out.Custom["ORIGINALYEAR"] = fmt.Sprintf("%d", release.ReleaseGroup.FirstReleaseDate.Year())
	out.Custom["RELEASESTATUS"] = strings.ToLower(release.Status)
	out.Custom["SCRIPT"] = release.TextRepresentation.Script

	// XXX library doesn't implement these yet -- guess I need to do a PR
	out.Composer = in.Composer

	// XXX library doesn't implement these yet -- guess I need to do a PR
	copyCustomsNotImpl := []string{"MusicBrainz Work Id", "WORK", "LYRICIST", "PRODUCER", "ENGINEER", "MIXER", "ISRC", "REMIXER", "WRITER"}
	for _, v := range copyCustomsNotImpl {
		if _, ok := in.Custom[v]; ok {
			out.Custom[v] = in.Custom[v]
		}
	}

	// XXX Process urls
	if x, ok := in.Custom["URL_DISCOGS_RELEASE_SITE"]; ok {
		out.Custom["URL_DISCOGS_RELEASE_SITE"] = x
	}

	return &out, true, nil
}
