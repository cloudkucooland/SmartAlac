package sa

import (
	"bufio"
	"fmt"
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
	metadata := mb5_query_query(c.mb5query, "release", "", searchQuery, 0, nil, nil)
	if metadata == nil {
		return nil, fmt.Errorf("search failed")
	}
	defer mb5_metadata_delete(metadata)

	releaseList := mb5_metadata_get_releaselist(metadata)
	if releaseList == nil {
		return nil, nil
	}

	count := mb5_release_list_size(releaseList)
	results := make([]MBRelease, 0, count)

	for i := 0; i < count; i++ {
		rel := mb5_release_list_item(releaseList, i)
		var r MBRelease

		r.ID = mb5String(mb5_release_get_id, unsafe.Pointer(rel))
		r.Title = mb5String(mb5_release_get_title, unsafe.Pointer(rel))
		r.Date = mb5String(mb5_release_get_date, unsafe.Pointer(rel))
		r.Country = mb5String(mb5_release_get_country, unsafe.Pointer(rel))
		r.Barcode = mb5String(mb5_release_get_barcode, unsafe.Pointer(rel))
		r.Disambiguation = mb5String(mb5_release_get_disambiguation, unsafe.Pointer(rel))

		ac := mb5_release_get_artistcredit(rel)
		r.Artist = c.fmtArtistCreditMB5(ac)

		ml := mb5_release_get_mediumlist(rel)
		if ml != nil {
			r.TrackCount = mb5_medium_list_get_trackcount(ml)
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
