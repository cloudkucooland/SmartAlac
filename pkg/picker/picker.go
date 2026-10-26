package picker

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sorrow446/go-mp4tag"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudkucooland/SmartAlac/pkg/sa"
	"github.com/irlndts/go-discogs"
)

type state int

const (
	stateInput state = iota
	stateSearch
	stateResults
	stateDiscogsResults
	stateTagging
	stateDone
)

type mode int

const (
	modeMusicBrainz mode = iota
	modeDiscogs
)

type model struct {
	state          state
	mode           mode
	curator        *sa.Curator
	targetDir      string
	toc            string
	discID         string
	isrcs          []string
	trackDurations []int // Durations in seconds parsed from TOC
	artistInput    textinput.Model
	albumInput     textinput.Model
	mcnInput       textinput.Model
	focusIndex     int
	results        []sa.MBRelease
	dgResults      []discogs.Result
	list           list.Model
	error          error
	status         string
	width          int
	height         int
}

// parseTOCDurations converts a libcdio TOC string into track durations in seconds.
func parseTOCDurations(tocStr string) []int {
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

func NewModel(c *sa.Curator, dir string) model {
	art := textinput.New()
	art.Placeholder = "Album Artist"
	art.Focus()

	alb := textinput.New()
	alb.Placeholder = "Album Title"

	mcn := textinput.New()
	mcn.Placeholder = "MCN / Barcode"

	var toc string
	var discID string
	var isrcs []string
	var durations []int

	ripdatapath := filepath.Join(dir, "ripdata.json")
	if data, err := os.ReadFile(ripdatapath); err == nil {
		var rd struct {
			TOC      string `json:"TOC"`
			MBdiscid string `json:"MBdiscid"`
			Tracks   []struct {
				ISRC string `json:"ISRC"`
			} `json:"Tracks"`
		}
		if json.Unmarshal(data, &rd) == nil {
			toc = rd.TOC
			discID = rd.MBdiscid
			for _, t := range rd.Tracks {
				if t.ISRC != "" {
					isrcs = append(isrcs, t.ISRC)
				}
			}
			if toc != "" {
				durations = parseTOCDurations(toc)
			}
		}
	}

	m := model{
		state:          stateInput,
		curator:        c,
		targetDir:      dir,
		toc:            toc,
		discID:         discID,
		isrcs:          isrcs,
		trackDurations: durations,
		artistInput:    art,
		albumInput:     alb,
		mcnInput:       mcn,
		focusIndex:     0,
	}

	files, _ := os.ReadDir(dir)
	var firstM4A string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".m4a") {
			firstM4A = filepath.Join(dir, f.Name())
			break
		}
	}

	if firstM4A != "" {
		mp4, err := mp4tag.Open(firstM4A)
		if err == nil {
			tags, err := mp4.Read()
			if err == nil {
				var mbid string
				if vals, ok := tags.Custom["MusicBrainz Album Id"]; ok && len(vals) > 0 {
					mbid = vals[0]
				}
				if m.toc == "" {
					if vals, ok := tags.Custom["TOC"]; ok && len(vals) > 0 {
						m.toc = vals[0]
						m.trackDurations = parseTOCDurations(m.toc)
					}
				}
				if m.discID == "" {
					if vals, ok := tags.Custom["MusicBrainz Disc Id"]; ok && len(vals) > 0 {
						m.discID = vals[0]
					}
				}
				var barcode string
				if vals, ok := tags.Custom["BARCODE"]; ok && len(vals) > 0 {
					barcode = vals[0]
				}
				if barcode == "" || sa.IsAllZeros(barcode) {
					if vals, ok := tags.Custom["MCN"]; ok && len(vals) > 0 {
						barcode = vals[0]
					}
				}
				if sa.IsAllZeros(barcode) {
					barcode = ""
				}
				artist := tags.AlbumArtist
				if artist == "" {
					artist = tags.Artist
				}
				album := tags.Album

				m.artistInput.SetValue(artist)
				m.albumInput.SetValue(album)
				m.mcnInput.SetValue(barcode)

				if mbid == "" && barcode != "" && m.toc == "" && m.discID == "" {
					m.mode = modeDiscogs
				}
			}
			mp4.Close()
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type searchResultMsg []sa.MBRelease
type discogsResultMsg []discogs.Result
type errorMsg error
type tagDoneMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+t":
			if m.state == stateInput {
				if m.mode == modeMusicBrainz {
					m.mode = modeDiscogs
				} else {
					m.mode = modeMusicBrainz
				}
			}
		}

		switch m.state {
		case stateInput:
			switch msg.String() {
			case "tab", "shift+tab", "up", "down":
				m.artistInput.Blur()
				m.albumInput.Blur()
				m.mcnInput.Blur()
				if msg.String() == "shift+tab" || msg.String() == "up" {
					m.focusIndex--
					if m.focusIndex < 0 {
						m.focusIndex = 2
					}
				} else {
					m.focusIndex++
					if m.focusIndex > 2 {
						m.focusIndex = 0
					}
				}
				switch m.focusIndex {
				case 0:
					m.artistInput.Focus()
				case 1:
					m.albumInput.Focus()
				case 2:
					m.mcnInput.Focus()
				}
			case "enter":
				m.state = stateSearch
				m.error = nil
				return m, func() tea.Msg {
					if m.mode == modeMusicBrainz {
						var res []sa.MBRelease
						var err error
						if m.toc != "" {
							res, err = m.curator.SearchByTOC(context.Background(), m.toc)
						}
						if (res == nil || len(res) == 0) && m.discID != "" {
							res, err = m.curator.SearchByTOC(context.Background(), m.discID)
						}
						if (res == nil || len(res) == 0) && len(m.isrcs) > 0 {
							res, err = m.curator.SearchByISRCs(context.Background(), m.isrcs)
						}
						if (res == nil || len(res) == 0) && m.toc == "" && m.discID == "" && len(m.isrcs) == 0 {
							res, err = m.curator.SearchMB(context.Background(), m.artistInput.Value(), m.albumInput.Value())
						}
						if err != nil {
							return errorMsg(err)
						}
						return searchResultMsg(res)
					} else {
						// Primary search by text / MCN fields
						res, err := m.curator.SearchDiscogs(context.Background(), m.artistInput.Value(), m.albumInput.Value(), m.mcnInput.Value())
						if (err != nil || len(res) == 0) && len(m.trackDurations) > 0 {
							// Fallback: search Discogs by Track Count and TOC duration array
							res, err = m.curator.SearchDiscogsByTrackDurations(context.Background(), m.trackDurations)
						}
						if err != nil {
							return errorMsg(err)
						}
						return discogsResultMsg(res)
					}
				}
			}
			var cmd tea.Cmd
			switch m.focusIndex {
			case 0:
				m.artistInput, cmd = m.artistInput.Update(msg)
			case 1:
				m.albumInput, cmd = m.albumInput.Update(msg)
			case 2:
				m.mcnInput, cmd = m.mcnInput.Update(msg)
			}
			return m, cmd

		case stateResults:
			if msg.String() == "enter" {
				selected := m.list.SelectedItem().(item).release
				m.state = stateTagging
				m.status = fmt.Sprintf("Tagging files in %s...", m.targetDir)
				return m, func() tea.Msg {
					err := m.curator.TagDirectory(context.Background(), m.targetDir, selected.ID, nil)
					if err != nil {
						return errorMsg(err)
					}
					return tagDoneMsg{}
				}
			}
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd

		case stateDiscogsResults:
			if msg.String() == "enter" {
				selected := m.list.SelectedItem().(dgItem).result
				m.state = stateTagging
				m.status = fmt.Sprintf("Tagging files in %s...", m.targetDir)
				return m, func() tea.Msg {
					files, err := os.ReadDir(m.targetDir)
					if err != nil {
						return errorMsg(err)
					}
					for _, f := range files {
						if strings.HasSuffix(f.Name(), ".m4a") {
							path := filepath.Join(m.targetDir, f.Name())
							mp4, err := mp4tag.Open(path)
							if err == nil {
								tags, _ := mp4.Read()
								if tags.TrackNumber < 1 {
									base := f.Name()
									if ext := filepath.Ext(base); ext != "" {
										base = strings.TrimSuffix(base, ext)
									}
									fields := strings.Fields(base)
									if len(fields) > 0 {
										if t, err := strconv.Atoi(fields[0]); err == nil && t > 0 {
											tags.TrackNumber = int16(t)
										}
									}
								}
								newTags, changed, _ := m.curator.UpdateFromDiscogs(context.Background(), tags, selected.ID)
								if changed {
									customKeys := make([]string, 0, len(newTags.Custom))
									for k := range newTags.Custom {
										customKeys = append(customKeys, k)
									}
									_ = mp4.Write(newTags, customKeys)
								}
								mp4.Close()
							}
						}
					}
					return tagDoneMsg{}
				}
			}
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

	case searchResultMsg:
		m.results = msg
		if len(m.results) == 0 {
			m.state = stateInput
			m.error = fmt.Errorf("no releases found")
			return m, nil
		}
		items := make([]list.Item, len(m.results))
		for i, r := range m.results {
			items[i] = item{r}
		}
		m.list = list.New(items, list.NewDefaultDelegate(), m.width, m.height-4)
		m.list.Title = "Select MusicBrainz Release"
		m.state = stateResults
		return m, nil

	case discogsResultMsg:
		m.dgResults = msg
		if len(m.dgResults) == 0 {
			m.state = stateInput
			m.error = fmt.Errorf("no releases found on Discogs")
			return m, nil
		}
		items := make([]list.Item, len(m.dgResults))
		for i, r := range m.dgResults {
			items[i] = dgItem{r}
		}
		m.list = list.New(items, list.NewDefaultDelegate(), m.width, m.height-4)
		m.list.Title = "Select Discogs Release"
		m.state = stateDiscogsResults
		return m, nil

	case tagDoneMsg:
		m.state = stateDone
		return m, tea.Quit

	case errorMsg:
		m.error = msg
		m.state = stateInput
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.state == stateResults || m.state == stateDiscogsResults {
			m.list.SetSize(msg.Width, msg.Height-4)
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	modeStr := "MusicBrainz"
	if m.mode == modeDiscogs {
		modeStr = "Discogs"
	}

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Bold(true).
		Render(fmt.Sprintf("PICKER: Album Tagger (%s)", modeStr))

	switch m.state {
	case stateInput:
		s = fmt.Sprintf(
			"Target Directory: %s\n\n%s\n%s\n%s\n\n(tab to switch, ctrl+t to toggle MB/DG, enter to search)\n",
			m.targetDir,
			m.artistInput.View(),
			m.albumInput.View(),
			m.mcnInput.View(),
		)

		if m.error != nil {
			s += lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Render(fmt.Sprintf("\nError: %v\n", m.error))
		}

	case stateSearch:
		s = fmt.Sprintf("\n  Searching %s...\n", modeStr)

	case stateResults, stateDiscogsResults:
		s = m.list.View()

	case stateTagging:
		s = fmt.Sprintf("\n  %s\n", m.status)

	case stateDone:
		s = "\n  Done! All files tagged.\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, "\n", s)
}

type item struct {
	release sa.MBRelease
}

func (i item) Title() string { return i.release.Title }
func (i item) Description() string {
	d := fmt.Sprintf("%s | %s | %d tracks", i.release.Artist, i.release.Date, i.release.TrackCount)
	if i.release.Label != "" {
		d += " | " + i.release.Label
	}
	if i.release.CatalogNumber != "" {
		d += " [" + i.release.CatalogNumber + "]"
	}
	if i.release.Media != "" {
		d += " (" + i.release.Media + ")"
	}
	if i.release.Disambiguation != "" {
		d += " — " + i.release.Disambiguation
	}
	return d
}
func (i item) FilterValue() string { return i.release.Title + " " + i.release.Artist }

type dgItem struct {
	result discogs.Result
}

func (i dgItem) Title() string { return i.result.Title }
func (i dgItem) Description() string {
	return fmt.Sprintf("ID: %d | Year: %s | Type: %s", i.result.ID, i.result.Year, i.result.Type)
}
func (i dgItem) FilterValue() string { return i.result.Title }
