package sa

import (
	"bufio"
	"fmt"
	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type MBRelease struct {
	ID             string
	Title          string
	Artist         string
	Date           string
	Country        string
	Barcode        string
	Disambiguation string
	TrackCount     int
	Label          string
	CatalogNumber  string
	Media          string
}

func (c *Curator) SearchMB(artist, album string) ([]MBRelease, error) {
	if c.mb5query == nil {
		return nil, fmt.Errorf("mb5 query not initialized")
	}

	// Lucene query format for MusicBrainz
	searchQuery := fmt.Sprintf("artist:\"%s\" AND release:\"%s\"", artist, album)
	if artist == "" {
		searchQuery = fmt.Sprintf("release:\"%s\"", album)
	}

	c.rl.Take()

	// Prepare query parameters: For search, use "query" parameter
	var params [1]*byte
	p1 := []byte("query")
	params[0] = &p1[0]

	var values [1]*byte
	v1 := []byte(searchQuery)
	values[0] = &v1[0]

	// Call QueryQuery with entity="release", id="", resource="", num_params=1
	metadata := mb5.QueryQuery(c.mb5query, "release", "", "", 1, unsafe.Pointer(&params), unsafe.Pointer(&values))
	if metadata == nil {
		return nil, fmt.Errorf("search failed (metadata nil)")
	}
	defer mb5.MetadataDelete(metadata)

	if result := mb5.QueryGetLastresult(c.mb5query); result != 0 {
		var errbuf [256]byte
		mb5.QueryGetLasterrormessage(c.mb5query, &errbuf[0], 256)
		return nil, fmt.Errorf("search error: %s", strings.Trim(string(errbuf[:]), "\x00"))
	}

	releaseList := mb5.MetadataGetReleaselist(metadata)
	if releaseList == nil {
		return nil, nil
	}

	count := mb5.ReleaseListSize(releaseList)
	results := make([]MBRelease, 0, count)

	for i := 0; i < count; i++ {
		rel := mb5.ReleaseListItem(releaseList, i)
		var r MBRelease

		r.ID = mb5.String(mb5.ReleaseGetID, unsafe.Pointer(rel))
		r.Title = mb5.String(mb5.ReleaseGetTitle, unsafe.Pointer(rel))
		r.Date = mb5.String(mb5.ReleaseGetDate, unsafe.Pointer(rel))
		r.Country = mb5.String(mb5.ReleaseGetCountry, unsafe.Pointer(rel))
		r.Barcode = mb5.String(mb5.ReleaseGetBarcode, unsafe.Pointer(rel))
		r.Disambiguation = mb5.String(mb5.ReleaseGetDisambiguation, unsafe.Pointer(rel))

		ac := mb5.ReleaseGetArtistcredit(rel)
		r.Artist = c.fmtArtistCreditMB5(ac)

		ml := mb5.ReleaseGetMediumlist(rel)
		if ml != nil {
			r.TrackCount = mb5.MediumListGetTrackcount(ml)
			// Get Media format from the first medium
			if mb5.MediumListSize(ml) > 0 {
				med := mb5.MediumListItem(ml, 0)
				r.Media = mb5.String(mb5.MediumGetFormat, unsafe.Pointer(med))
			}
		}

		liList := mb5.ReleaseGetLabelinfolist(rel)
		if liList != nil {
			r.Label = c.fmtLabelsMB5(liList)
			r.CatalogNumber = c.fmtCatalogNumbersMB5(liList)
		}

		results = append(results, r)
	}

	return results, nil
}

func (c *Curator) SelectRelease(results []MBRelease) (*MBRelease, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	fmt.Println("\nMusicBrainz Search Results:")
	for i, r := range results {
		fmt.Printf("[%d] %s - %s\n", i+1, r.Artist, r.Title)
		fmt.Printf("    ID: %s | Date: %s | Country: %s | Tracks: %d\n", r.ID, r.Date, r.Country, r.TrackCount)
		if r.Label != "" || r.CatalogNumber != "" {
			fmt.Printf("    Label: %s | Cat#: %s\n", r.Label, r.CatalogNumber)
		}
		if r.Media != "" {
			fmt.Printf("    Media: %s\n", r.Media)
		}
		if r.Disambiguation != "" {
			fmt.Printf("    Disambiguation: %s\n", r.Disambiguation)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nSelect a release (1-%d, or 's' to skip): ", len(results))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "s" {
			return nil, nil
		}

		idx, err := strconv.Atoi(input)
		if err == nil && idx >= 1 && idx <= len(results) {
			return &results[idx-1], nil
		}

		fmt.Println("Invalid selection, please try again.")
	}
}
