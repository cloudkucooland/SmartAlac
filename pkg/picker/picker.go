package picker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	state       state
	mode        mode
	curator     *sa.Curator
	targetDir   string
	artistInput textinput.Model
	albumInput  textinput.Model
	mcnInput    textinput.Model
	focusIndex  int
	results     []sa.MBRelease
	dgResults   []discogs.Result
	list        list.Model
	error       error
	status      string
	width       int
	height      int
}

func NewModel(c *sa.Curator, dir string) model {
	art := textinput.New()
	art.Placeholder = "Album Artist"
	art.Focus()

	alb := textinput.New()
	alb.Placeholder = "Album Title"

	mcn := textinput.New()
	mcn.Placeholder = "MCN / Barcode"

	m := model{
		state:       stateInput,
		curator:     c,
		targetDir:   dir,
		artistInput: art,
		albumInput:  alb,
		mcnInput:    mcn,
		focusIndex:  0,
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
				mbid := tags.Custom["MusicBrainz Album Id"]
				barcode := tags.Custom["BARCODE"]
				if barcode == "" || sa.IsAllZeros(barcode) {
					barcode = tags.Custom["MCN"]
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

				if mbid == "" && barcode != "" {
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
						res, err := m.curator.SearchMB(context.Background(), m.artistInput.Value(), m.albumInput.Value())
						if err != nil {
							return errorMsg(err)
						}
						return searchResultMsg(res)
					} else {
						res, err := m.curator.SearchDiscogs(context.Background(), m.artistInput.Value(), m.albumInput.Value(), m.mcnInput.Value())
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
					// TagDirectory will now handle Discogs fallback internally if MBID is empty
					// but since we selected a specific Discogs ID, we should pass it.
					// We'll need a way to tell TagDirectory to use a specific Discogs ID.
					// For now, let's update all files in dir with this Discogs ID.
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
		if m.state == stateResults {
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

	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Bold(true).Render(fmt.Sprintf("PICKER: Album Tagger (%s)", modeStr))

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
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(fmt.Sprintf("\nError: %v\n", m.error))
		}

	case stateSearch:
		s = fmt.Sprintf("\n  Searching %s...\n", modeStr)

	case stateResults, stateDiscogsResults:
		s = m.list.View()

	case stateTagging:
		s = fmt.Sprintf("\n  %s\n", m.status)

	case stateDone:
		s = "\n  Done! All files tagged and moved.\n"
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
