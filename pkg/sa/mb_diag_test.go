package sa

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
)

func TestMusicBrainzDataDump(t *testing.T) {
	// Only run if specifically requested, as it needs the shared library
	// go test -v -run TestMusicBrainzDataDump
	
	err := mb5.Init()
	if err != nil {
		t.Skip("libmusicbrainz5 not found, skipping live test")
	}

	query := mb5.QueryNew("SmartAlac-Diag", "musicbrainz.org", 0)
	if query == nil {
		t.Fatal("failed to create query object")
	}
	defer mb5.QueryDelete(query)

	// Tom Lehrer - That Was The Year That Was (Mono Vinyl)
	releaseID := "a431d9dc-b7e4-48db-a06f-ddccf2f26d81"

	var params [1]*byte
	p1 := []byte("inc")
	params[0] = &p1[0]

	var values [1]*byte
	v1 := []byte("artists labels recordings release-groups artist-credits")
	values[0] = &v1[0]

	metadata := mb5.QueryQuery(query, "release", releaseID, "", 1, unsafe.Pointer(&params), unsafe.Pointer(&values))
	if metadata == nil {
		t.Fatal("query failed")
	}
	defer mb5.MetadataDelete(metadata)

	release := mb5.MetadataGetRelease(metadata)
	if release == nil {
		t.Fatal("no release in metadata")
	}

	title := mb5.String(mb5.ReleaseGetTitle, unsafe.Pointer(release))
	fmt.Printf("Release Title: %q\n", title)

	ac := mb5.ReleaseGetArtistcredit(release)
	fmt.Printf("Artist Credit Pointer: %p\n", ac)

	if ac != nil {
		ncl := mb5.ArtistcreditGetNamecreditlist(ac)
		count := mb5.NamecreditListSize(ncl)
		fmt.Printf("Name Credit List Count: %d\n", count)
		for i := 0; i < count; i++ {
			nc := mb5.NamecreditListItem(ncl, i)
			name := mb5.String(mb5.NamecreditGetName, unsafe.Pointer(nc))
			join := mb5.String(mb5.NamecreditGetJoinphrase, unsafe.Pointer(nc))
			fmt.Printf("  Credit %d: Name=%q, Join=%q\n", i, name, join)
			
			artist := mb5.NamecreditGetArtist(nc)
			if artist != nil {
				aname := mb5.String(mb5.ArtistGetName, unsafe.Pointer(artist))
				fmt.Printf("    Artist Name: %q\n", aname)
			}
		}
	} else {
		fmt.Println("CRITICAL: Artist Credit is NIL for the release")
	}
}
