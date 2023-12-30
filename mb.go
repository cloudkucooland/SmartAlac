package sa

import (
	"github.com/Sorrow446/go-mp4tag"
	"github.com/kr/pretty"
	"github.com/michiwend/gomusicbrainz"
	"log"
	// if the caches are too big, try one of these (or others)
	// git.mills.io/prologic/bitcask
	// github.com/dgraph-io/ristretto
	// github.com/peterbourgon/diskv
)

/* go-mp4tag.MP4Tag:
   Album:What Girls Want
   AlbumSort:
   AlbumArtist:Material Issue
   AlbumArtistSort:
   Artist:Material Issue Artist
   Sort:
   BPM:-1
   Comment:alternative rock
   Composer:
   ComposerSort:
   Conductor:
   Copyright:
   CustomGenre:Alternative Rock
   Date:
   Description:
   Director:
   DiscNumber:1
   DiscTotal:1
   Genre:0
   ItunesAdvisory:0
   ItunesAlbumID:-1
   ItunesArtistID:-1
   Lyrics:
   Narrator:
   Pictures:[]
   Publisher:
   Title:What Girls Want (radio mix)
   TitleSort:
   TrackNumber:1
   TrackTotal:3
   Year:1992

   --- Custom ---
   ARTISTS:Material Issue
   BARCODE:314512333243
   CATALOGNUMBER:CDP 685
   Country:United States
   KEY:A
   LABEL:Mercury Records
   LANGUAGE:eng
   MEDIA:CD
   MOOD:bright; happy; relaxed
   MusicBrainz Album Artist Id:35290bf2-aa8e-412d-8f4e-0b3362044958
   MusicBrainz Album Id:bf7cb22d-a15a-4ce2-a5ff-b0c4f8608f2e
   MusicBrainz Album Release Country:US
   MusicBrainz Album Type:single
   MusicBrainz Artist Id:35290bf2-aa8e-412d-8f4e-0b3362044958
   MusicBrainz Disc Id:YeZdB4dY7lQIC0ZSgkbe2R4HXAs-
   MusicBrainz Release Group Id:c2912250-8518-45eb-bf43-5d57d8607d82
   MusicBrainz Release Track Id:5bbd0145-a8d8-3321-8442-638690620ae7
   MusicBrainz Track Id:5bbd0145-a8d8-3321-8442-638690620ae7
   ORIGINALDATE:1992
   ORIGINALYEAR:1992
   RELEASESTATUS:promotion
   SCRIPT:Latn
   URL_DISCOGS_ARTIST_SITE:http://discogs.com/artist/261512
   URL_DISCOGS_RELEASE_SITE:https://www.discogs.com/release/2653425
   URL_LYRICS_SITE:http://lyricsfly.com/search/correction.php?8ebcad6785&id=621629
*/

var client *gomusicbrainz.WS2Client

var releases map[string]*gomusicbrainz.Release

// var artists map[string]*gomusicbrainz.Artist
// var labels map[string]*gomusicbrainz.Label

func init() {
	client, _ = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"SmartAlac",
		"0.0.0",
		"http://github.com/cloudkucooland/SmartAlac")

	releases = make(map[string]*gomusicbrainz.Release)
	// artists = make(map[string]*gomusicbrainz.Artist)
	// labels = make(map[string]*gomusicbrainz.Label)
}

func updateFromMB(in *mp4tag.MP4Tags) (*mp4tag.MP4Tags, bool, error) {
	out := mp4tag.MP4Tags{}
	// constants
	out.ItunesAdvisory = 0
	out.ItunesAlbumID = -1
	out.ItunesArtistID = -1

	// things which we don't touch but do want to preserve
	out.BPM = in.BPM
	out.Lyrics = in.Lyrics

	// copy over custom items determined eleswhere
	out.Custom = make(map[string]string)
	copyCustoms := []string{"KEY", "MOOD"}
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
		log.Printf("querying for %s\n\n", releaseid)
		var err error
		release, err = client.LookupRelease(gomusicbrainz.MBID(releaseid), "artist-credits", "recordings", "release-groups", "media", "isrcs", "release-rels", "release-group-rels", "url-rels")
		if err != nil {
			log.Println(err.Error())
			return in, false, err
		}
		releases[releaseid] = release
		// log.Printf("%# v\n", pretty.Formatter(release))
	}

	// per release items
	out.Album = release.Title
	// out.AlbumSort =
	out.AlbumArtist = fmtArtistCredit(release.ArtistCredit.NameCredits)
	// out.AlbumArtistSort =
	out.DiscNumber = in.DiscNumber
	// out.DiscTotal = release.Mediums.Length()
	// out.Sort:
	out.TrackTotal = getMediumTrackCount(release, in.DiscNumber)
	// out.Publisher = release.
	out.Year = int32(release.ReleaseGroup.FirstReleaseDate.Year())

	// per track items
	track := getTrack(release, gomusicbrainz.MBID(tid))
	out.Artist = fmtArtistCredit(track.Recording.ArtistCredit.NameCredits)
	/*   out.Comment:alternative rock
	out.Composer:
	out.ComposerSort:
	out.Conductor:
	out.Copyright:
	out.CustomGenre:Alternative Rock
	out.Date:
	out.Description:
	out.Director:
	out.Genre:0
	out.Lyrics:
	out.Narrator:
	out.Title:What Girls Want (radio mix)
	out.TitleSort:
	*/
	out.TrackNumber = in.TrackNumber

	log.Printf("%# v\n", pretty.Formatter(pretty.Diff(in, &out)))

	return &out, true, nil
}

func getMediumTrackCount(r *gomusicbrainz.Release, discnumber int16) int16 {
	return 0
}

func fmtArtistCredit(a []gomusicbrainz.NameCredit) string {
	var s string
	for _, v := range a {
		s += v.Artist.Name
	}
	log.Println(s)
	return s
}

func fmtArtistCreditSort(a []gomusicbrainz.NameCredit) string {
	var s string
	for _, v := range a {
		s += v.Artist.SortName
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
