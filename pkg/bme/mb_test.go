package bme

import (
	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
	"strings"
	"testing"
	"unsafe"
)

func TestMBLookup(t *testing.T) {
	// Mock implementation for testing
	dummy := 1
	mb5.QueryNew = func(a, b string, c int) mb5.Query { return mb5.Query(unsafe.Pointer(&dummy)) }
	mb5.QueryQuery = func(a mb5.Query, b, c, d string, e int, f, g unsafe.Pointer) mb5.Metadata {
		return mb5.Metadata(unsafe.Pointer(&dummy))
	}
	mb5.MetadataClone = func(a mb5.Metadata) mb5.Metadata { return a }
	mb5.MetadataDelete = func(a mb5.Metadata) {}
	mb5.QueryGetLastresult = func(a mb5.Query) mb5.QueryResult { return 0 }
	mb5.MetadataGetDisc = func(a mb5.Metadata) mb5.Disc { return mb5.Disc(unsafe.Pointer(&dummy)) }
	mb5.DiscClone = func(a mb5.Disc) mb5.Disc { return a }
	mb5.DiscDelete = func(a mb5.Disc) {}
	mb5.DiscGetReleaselist = func(a mb5.Disc) mb5.ReleaseList { return mb5.ReleaseList(unsafe.Pointer(&dummy)) }
	mb5.ReleaseListClone = func(a mb5.ReleaseList) mb5.ReleaseList { return a }
	mb5.ReleaseListDelete = func(a mb5.ReleaseList) {}
	mb5.ReleaseListSize = func(a mb5.ReleaseList) int { return 1 }
	mb5.ReleaseListItem = func(a mb5.ReleaseList, b int) mb5.Release { return mb5.Release(unsafe.Pointer(&dummy)) }
	mb5.ReleaseClone = func(a mb5.Release) mb5.Release { return a }
	mb5.ReleaseDelete = func(a mb5.Release) {}
	mb5.ReleaseGetID = func(a unsafe.Pointer, b *byte, c int) int {
		copy((*[37]byte)(unsafe.Pointer(b))[:], "test-release-id")
		return 15
	}
	mb5.ReleaseGetTitle = func(a unsafe.Pointer, b *byte, c int) int {
		copy((*[256]byte)(unsafe.Pointer(b))[:], "Test Album")
		return 10
	}
	mb5.MetadataGetRelease = func(a mb5.Metadata) mb5.Release { return mb5.Release(unsafe.Pointer(&dummy)) }
	mb5.ReleaseGetCountry = func(a unsafe.Pointer, b *byte, c int) int { return 0 }
	mb5.ReleaseGetBarcode = func(a unsafe.Pointer, b *byte, c int) int { return 0 }
	mb5.ReleaseGetDisambiguation = func(a unsafe.Pointer, b *byte, c int) int { return 0 }
	mb5.ReleaseMediaMatchingDiscid = func(a mb5.Release, b string) mb5.MediumList {
		return mb5.MediumList(unsafe.Pointer(&dummy))
	}
	mb5.MediumListSize = func(a mb5.MediumList) int { return 1 }
	mb5.MediumListItem = func(a mb5.MediumList, b int) mb5.Medium { return mb5.Medium(unsafe.Pointer(&dummy)) }
	mb5.MediumGetTracklist = func(a mb5.Medium) mb5.TrackList { return mb5.TrackList(unsafe.Pointer(&dummy)) }
	mb5.TrackListClone = func(a mb5.TrackList) mb5.TrackList { return a }
	mb5.TrackListDelete = func(a mb5.TrackList) {}
	mb5.TrackListGetCount = func(a mb5.TrackList) int { return 1 }
	mb5.MediumGetPosition = func(a mb5.Medium) int { return 1 }
	mb5.TrackListItem = func(a mb5.TrackList, b int) mb5.Track { return mb5.Track(unsafe.Pointer(&dummy)) }
	mb5.TrackGetPosition = func(a mb5.Track) int { return 1 }
	mb5.TrackGetRecording = func(a mb5.Track) mb5.Recording { return mb5.Recording(unsafe.Pointer(&dummy)) }
	mb5.RecordingGetID = func(a unsafe.Pointer, b *byte, c int) int { return 0 }
	mb5.RecordingGetTitle = func(a unsafe.Pointer, b *byte, c int) int {
		copy((*[256]byte)(unsafe.Pointer(b))[:], "Test Track")
		return 10
	}
	mb5.RecordingGetArtistcredit = func(a mb5.Recording) mb5.ArtistCredit {
		return mb5.ArtistCredit(unsafe.Pointer(&dummy))
	}
	mb5.ArtistcreditGetNamecreditlist = func(a mb5.ArtistCredit) mb5.NameCreditList {
		return mb5.NameCreditList(unsafe.Pointer(&dummy))
	}
	mb5.NamecreditListGetCount = func(a mb5.NameCreditList) int {
		return 2
	}

	names := []string{"Artist A", "Artist B"}
	joins := []string{" & ", ""}
	mb5.NamecreditListItem = func(a mb5.NameCreditList, b int) mb5.NameCredit {
		return mb5.NameCredit(uintptr(b + 1)) // Shift by 1 to avoid 0 (nil)
	}
	mb5.NamecreditGetName = func(a unsafe.Pointer, b *byte, c int) int {
		idx := int(uintptr(a)) - 1
		buf := (*[256]byte)(unsafe.Pointer(b))
		for i := range buf {
			buf[i] = 0
		}
		if idx >= 0 && idx < len(names) {
			copy(buf[:], names[idx])
			return len(names[idx])
		}
		return 0
	}
	mb5.NamecreditGetJoinphrase = func(a unsafe.Pointer, b *byte, c int) int {
		idx := int(uintptr(a)) - 1
		buf := (*[256]byte)(unsafe.Pointer(b))
		for i := range buf {
			buf[i] = 0
		}
		if idx >= 0 && idx < len(joins) {
			copy(buf[:], joins[idx])
			return len(joins[idx])
		}
		return 0
	}
	mb5.NamecreditGetArtist = func(a mb5.NameCredit) mb5.Artist { return mb5.Artist(unsafe.Pointer(&dummy)) }
	mb5.ArtistGetName = func(a unsafe.Pointer, b *byte, c int) int {
		buf := (*[256]byte)(unsafe.Pointer(b))
		for i := range buf {
			buf[i] = 0
		}
		return 0
	}

	releases := mb_lookup_discid("test-discid", 1)

	if len(releases) != 1 {
		t.Fatalf("expected 1 release, got %d", len(releases))
	}

	mbr := releases[0]
	if mbr.Title != "Test Album" {
		t.Errorf("expected title Test Album, got %s", mbr.Title)
	}

	if len(mbr.Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(mbr.Tracks))
	}

	expectedArtist := "Artist A & Artist B"
	if strings.TrimSpace(mbr.Tracks[0].Artist) != expectedArtist {
		t.Errorf("expected artist %s, got '%s'", expectedArtist, mbr.Tracks[0].Artist)
	}
}
