package bme

import (
	"log/slog"
	"strings"
	"unsafe"

	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
)

type mb_release struct {
	DiscID         string
	ReleaseID      string
	AlbumArtist    string
	Title          string
	DiscPosition   int
	Tracks         []mb_track
	Country        string
	Barcode        string
	Disambiguation string
}

type mb_track struct {
	Position int
	TrackID  string
	Artist   string
	Title    string
}

func mb_lookup_discid(mbid string, expectedTracks int) []mb_release {
	var releases []mb_release

	query := mb5.QueryNew("bme-tag-0.0", "musicbrainz.org", 0)
	if query == nil {
		slog.Error("mb_lookup_discid: unable to get query")
		return releases
	}

	metadata1 := mb5.QueryQuery(query, "discid", mbid, "", 0, nil, nil)
	if metadata1 == nil {
		slog.Debug("mb_lookup_discid", "msg", "no results")
		return releases
	}
	metadata1 = mb5.MetadataClone(metadata1)
	defer mb5.MetadataDelete(metadata1)

	result := mb5.QueryGetLastresult(query)
	if result != 0 {
		mb_error_message("last query result", query)
		return releases
	}

	disc := mb5.MetadataGetDisc(metadata1)
	if disc == nil {
		mb_error_message("get_disc", query)
		return releases
	}

	disc = mb5.DiscClone(disc)
	defer mb5.DiscDelete(disc)

	rl := mb5.DiscGetReleaselist(disc)
	if rl == nil {
		mb_error_message("get_releaselist", query)
		return releases
	}
	rl = mb5.ReleaseListClone(rl)
	defer mb5.ReleaseListDelete(rl)

	rcount := mb5.ReleaseListSize(rl)
	for e := 0; e < rcount; e++ {
		var mbr mb_release
		mbr.DiscID = mbid

		shortrelease := mb5.ReleaseListItem(rl, e)
		shortrelease = mb5.ReleaseClone(shortrelease)
		defer mb5.ReleaseDelete(shortrelease)

		var releaseID [37]byte
		mb5.ReleaseGetID(unsafe.Pointer(shortrelease), (*byte)(unsafe.Pointer(&releaseID[0])), 37)
		mbr.ReleaseID = strings.Trim(string(releaseID[:]), "\x00")

		var title [256]byte
		mb5.ReleaseGetTitle(unsafe.Pointer(shortrelease), (*byte)(unsafe.Pointer(&title[0])), 256)
		mbr.Title = strings.Trim(string(title[:]), "\x00")

		var params [1]*byte
		p1 := []byte("inc")
		params[0] = &p1[0]

		var values [1]*byte
		v1 := []byte("artists labels recordings release-groups url-rels discids artist-credits")
		values[0] = &v1[0]

		metadata2 := mb5.QueryQuery(query, "release", mbr.ReleaseID, "", 1, unsafe.Pointer(&params), unsafe.Pointer(&values))
		if metadata2 == nil {
			mb_error_message("metadata2 nil", query)
			continue
		}
		metadata2 = mb5.MetadataClone(metadata2)
		defer mb5.MetadataDelete(metadata2)

		fullrelease := mb5.MetadataGetRelease(metadata2)
		if fullrelease == nil {
			mb_error_message("full release nil", query)
			continue
		}

		if mbr.Title == "" {
			mb5.ReleaseGetTitle(unsafe.Pointer(fullrelease), (*byte)(unsafe.Pointer(&title[0])), 256)
			mbr.Title = strings.Trim(string(title[:]), "\x00")
		}

		var country [10]byte
		mb5.ReleaseGetCountry(unsafe.Pointer(fullrelease), (*byte)(unsafe.Pointer(&country[0])), 10)
		mbr.Country = strings.Trim(string(country[:]), "\x00")

		var barcode [16]byte
		mb5.ReleaseGetBarcode(unsafe.Pointer(fullrelease), (*byte)(unsafe.Pointer(&barcode[0])), 10)
		mbr.Barcode = strings.Trim(string(barcode[:]), "\x00")

		var disambiguation [256]byte
		mb5.ReleaseGetDisambiguation(unsafe.Pointer(fullrelease), (*byte)(unsafe.Pointer(&disambiguation[0])), 10)
		mbr.Disambiguation = strings.Trim(string(disambiguation[:]), "\x00")

		medialist := mb5.ReleaseMediaMatchingDiscid(fullrelease, mbid)
		if medialist == nil {
			mb_error_message("medialist nil", query)
			continue
		}

		mls := mb5.MediumListSize(medialist)
		if mls == 0 {
			mb_error_message("zero medialist items", query)
			continue
		}

		medium := mb5.MediumListItem(medialist, 0)
		if medium == nil {
			mb_error_message("medium nil", query)
			continue
		}

		tracklist := mb5.MediumGetTracklist(medium)
		if tracklist == nil {
			continue
		}
		tracklist = mb5.TrackListClone(tracklist)
		defer mb5.TrackListDelete(tracklist)

		trackcount := mb5.TrackListGetCount(tracklist)
		if expectedTracks > 0 && trackcount != expectedTracks {
			slog.Debug("skipping release due to track count mismatch", "title", mbr.Title, "mb_tracks", trackcount, "cd_tracks", expectedTracks)
			continue
		}

		mbr.DiscPosition = mb5.MediumGetPosition(medium)

		for j := 0; j < trackcount; j++ {
			var tmp mb_track

			track := mb5.TrackListItem(tracklist, j)
			if track == nil {
				continue
			}

			tmp.Position = mb5.TrackGetPosition(track)

			rec := mb5.TrackGetRecording(track)
			if rec != nil {
				var buf [256]byte
				mb5.RecordingGetID(unsafe.Pointer(rec), (*byte)(unsafe.Pointer(&buf[0])), 255)
				tmp.TrackID = strings.Trim(string(buf[:]), "\x00")

				var title [256]byte
				mb5.RecordingGetTitle(unsafe.Pointer(rec), (*byte)(unsafe.Pointer(&title[0])), 255)
				tmp.Title = strings.Trim(string(title[:]), "\x00")
			} else {
				var title [256]byte
				mb5.TrackGetTitle(unsafe.Pointer(track), (*byte)(unsafe.Pointer(&title[0])), 255)
				tmp.Title = strings.Trim(string(title[:]), "\x00")
			}

			ac := mb5.RecordingGetArtistcredit(rec)
			if ac == nil {
				ac = mb5.TrackGetArtistcredit(track)
			}

			if ac != nil {
				var fullartistname strings.Builder
				ncl := mb5.ArtistcreditGetNamecreditlist(ac)
				credits := mb5.NamecreditListGetCount(ncl)

				for k := 0; k < credits; k++ {
					nc := mb5.NamecreditListItem(ncl, k)
					if nc == nil {
						continue
					}

					var buf [256]byte
					mb5.NamecreditGetName(unsafe.Pointer(nc), (*byte)(unsafe.Pointer(&buf[0])), 256)
					n := strings.Trim(string(buf[:]), "\x00")
					if n != "" {
						fullartistname.WriteString(n)
					}
					mb5.NamecreditGetJoinphrase(unsafe.Pointer(nc), (*byte)(unsafe.Pointer(&buf[0])), 256)
					n = strings.Trim(string(buf[:]), "\x00")
					if n != "" {
						fullartistname.WriteString(n)
					} else {
						artist := mb5.NamecreditGetArtist(nc)
						mb5.ArtistGetName(unsafe.Pointer(artist), (*byte)(unsafe.Pointer(&buf[0])), 256)
						n = strings.Trim(string(buf[:]), "\x00")
						fullartistname.WriteString(n)
					}
				}
				tmp.Artist = fullartistname.String()
			}
			mbr.Tracks = append(mbr.Tracks, tmp)
		}
		releases = append(releases, mbr)
	}

	return releases
}

func mb_error_message(msg string, query mb5.Query) {
	result := mb5.QueryGetLastresult(query)
	slog.Debug("mb_lookup_discid", "last query result", result)

	var errbuf [256]byte
	mb5.QueryGetLasterrormessage(query, &errbuf[0], 256)
	slog.Debug("mb_lookup_discid", "msg", msg, "err", strings.Trim(string(errbuf[:]), "\x00"))
}
