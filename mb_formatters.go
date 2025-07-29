package sa

import (
	"log"
	"sort"
	"strings"

	"github.com/michiwend/gomusicbrainz"
)

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

// XXX library doesn't give connecting strings...
func fmtArtistCredit(a []gomusicbrainz.NameCredit) string {
	var s string
	for k, v := range a {
		if k > 0 {
			s += ", "
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

func getTrackRecordingID(r *gomusicbrainz.Release, trackID gomusicbrainz.MBID) gomusicbrainz.MBID {
	t := getTrack(r, trackID)
	if t == nil {
		return ""
	}
	return t.Recording.ID
}

func fmtCatalogNumber(l []gomusicbrainz.LabelInfo) string {
	// fast-path for the normal case
	if len(l) == 1 {
		return l[0].CatalogNumber
	}

	var s string

	reduce := make(map[string]bool)
	for _, li := range l {
		reduce[li.CatalogNumber] = true
	}
	m := make([]string, 0)
	for n := range reduce {
		m = append(m, n)
	}
	sort.Strings(m)

	for _, number := range m {
		if len(s) > 0 {
			s += "; "
		}
		s += number
	}
	return s
}

func fmtLabel(l []gomusicbrainz.LabelInfo) string {
	// fast-path for the normal case
	if len(l) == 1 && l[0].Label != nil && l[0].Label.Name != "" {
		return l[0].Label.Name
	}

	var s string

	reduce := make(map[string]bool)
	for _, li := range l {
        if li.Label != nil && li.Label.Name != "" {
		    reduce[li.Label.Name] = true
        }
	}
	m := make([]string, 0)
	for n := range reduce {
		m = append(m, n)
	}
	sort.Strings(m)

	for _, name := range m {
		if len(s) > 0 {
			s += "; "
		}
		s += name
	}
	return s
}

func mediumFormat(f string) string {
	return strings.Replace(f, "\"", "″", 1)
}

func formatDate(d gomusicbrainz.BrainzTime) string {
	switch d.Accuracy {
	case gomusicbrainz.Year:
		return d.Format("2006")
	case gomusicbrainz.Month:
		return d.Format("2006-01")
	case gomusicbrainz.Day:
		return d.Format("2006-01-02")
	}
	return d.Format("2006")
}
