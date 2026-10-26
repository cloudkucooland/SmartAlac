package sa

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

// ParseTOCDurations converts a libcdio TOC string into track durations in seconds.
func ParseTOCDurations(tocStr string) []int {
	parts := strings.Fields(tocStr)
	if len(parts) < 4 {
		return nil
	}

	leadout, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	var offsets []int
	for _, p := range parts[3:] {
		if val, err := strconv.Atoi(p); err == nil {
			offsets = append(offsets, val)
		}
	}

	if len(offsets) == 0 {
		return nil
	}

	durations := make([]int, len(offsets))
	for i := 0; i < len(offsets); i++ {
		var nextOffset int
		if i+1 < len(offsets) {
			nextOffset = offsets[i+1]
		} else {
			nextOffset = leadout
		}
		frames := nextOffset - offsets[i]
		// 75 CDDA frames per second
		durations[i] = int(math.Round(float64(frames) / 75.0))
	}
	return durations
}

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

func (c *Curator) SearchMB(ctx context.Context, artist, album string) ([]MBRelease, error) {
	if c.mb5query == nil {
		return nil, fmt.Errorf("mb5 query not initialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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
		r := c.parseReleaseMB5(rel)
		results = append(results, r)
	}

	return results, nil
}

func (c *Curator) parseReleaseMB5(rel mb5.Release) MBRelease {
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
	return r
}

func (c *Curator) SearchByISRCs(ctx context.Context, isrcs []string) ([]MBRelease, error) {
	if c.mb5query == nil {
		return nil, fmt.Errorf("mb5 query not initialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	type releaseScore struct {
		release MBRelease
		score   int
	}
	releaseMap := make(map[string]*releaseScore)

	for _, isrc := range isrcs {
		isrc = strings.TrimSpace(isrc)
		if isrc == "" {
			continue
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		searchQuery := fmt.Sprintf("isrc:%s", isrc)
		c.rl.Take()

		var params [1]*byte
		p1 := []byte("query")
		params[0] = &p1[0]

		var values [1]*byte
		v1 := []byte(searchQuery)
		values[0] = &v1[0]

		metadata := mb5.QueryQuery(c.mb5query, "recording", "", "", 1, unsafe.Pointer(&params), unsafe.Pointer(&values))
		if metadata == nil {
			continue
		}

		if result := mb5.QueryGetLastresult(c.mb5query); result != 0 {
			mb5.MetadataDelete(metadata)
			continue
		}

		recList := mb5.MetadataGetRecordinglist(metadata)
		if recList == nil {
			mb5.MetadataDelete(metadata)
			continue
		}

		recCount := mb5.RecordingListSize(recList)
		for i := 0; i < recCount; i++ {
			rec := mb5.RecordingListItem(recList, i)
			relList := mb5.RecordingGetReleaselist(rec)
			if relList == nil {
				continue
			}

			relCount := mb5.ReleaseListSize(relList)
			for j := 0; j < relCount; j++ {
				rel := mb5.ReleaseListItem(relList, j)
				parsed := c.parseReleaseMB5(rel)
				if parsed.ID == "" {
					continue
				}

				if rs, exists := releaseMap[parsed.ID]; exists {
					rs.score++
				} else {
					releaseMap[parsed.ID] = &releaseScore{
						release: parsed,
						score:   1,
					}
				}
			}
		}
		mb5.MetadataDelete(metadata)
	}

	if len(releaseMap) == 0 {
		return nil, nil
	}

	results := make([]MBRelease, 0, len(releaseMap))
	for _, rs := range releaseMap {
		results = append(results, rs.release)
	}

	sort.Slice(results, func(i, j int) bool {
		rsI := releaseMap[results[i].ID]
		rsJ := releaseMap[results[j].ID]
		if rsI.score != rsJ.score {
			return rsI.score > rsJ.score
		}
		diffI := abs(results[i].TrackCount - len(isrcs))
		diffJ := abs(results[j].TrackCount - len(isrcs))
		return diffI < diffJ
	})

	return results, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (c *Curator) SelectRelease(ctx context.Context, results []MBRelease) (*MBRelease, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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

func (c *Curator) SearchMBQuery(ctx context.Context, searchQuery string) ([]MBRelease, error) {
	if c.mb5query == nil {
		return nil, fmt.Errorf("mb5 query not initialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c.rl.Take()

	var params [1]*byte
	p1 := []byte("query")
	params[0] = &p1[0]

	var values [1]*byte
	v1 := []byte(searchQuery)
	values[0] = &v1[0]

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
		r := c.parseReleaseMB5(rel)
		results = append(results, r)
	}

	return results, nil
}

// SearchMBByDurations constructs a MusicBrainz Lucene query based on track count and durations.
func (c *Curator) SearchMBByDurations(ctx context.Context, durations []int) ([]MBRelease, error) {
	if len(durations) == 0 {
		return nil, fmt.Errorf("no durations provided")
	}

	// Build Lucene query: tracks:N AND dur:MS1 AND dur:MS2...
	// Allow a +/- 2000ms (2s) window per track
	var durQueries []string
	// Use first 5 tracks to avoid query-length limits on large track counts
	maxCheck := 5
	if len(durations) < maxCheck {
		maxCheck = len(durations)
	}

	for i := 0; i < maxCheck; i++ {
		ms := durations[i] * 1000
		durQueries = append(durQueries, fmt.Sprintf("dur:[%d TO %d]", ms-2000, ms+2000))
	}

	query := fmt.Sprintf("tracks:%d AND %s", len(durations), strings.Join(durQueries, " AND "))

	// Query MusicBrainz Release endpoint with structured query
	return c.SearchMBQuery(ctx, query)
}
