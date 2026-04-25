package picker

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudkucooland/SmartAlac/pkg/sa"
)

type state int

const (
	stateInput state = iota
	stateSearch
	stateResults
	stateTagging
	stateDone
)

type model struct {
	state       state
	curator     *sa.Curator
	targetDir   string
	artistInput textinput.Model
	albumInput  textinput.Model
	focusIndex  int
	results     []sa.MBRelease
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

	return model{
		state:       stateInput,
		curator:     c,
		targetDir:   dir,
		artistInput: art,
		albumInput:  alb,
		focusIndex:  0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type searchResultMsg []sa.MBRelease
type errorMsg error
type tagDoneMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.state {
		case stateInput:
			switch msg.String() {
			case "tab", "shift+tab", "up", "down":
				if m.focusIndex == 0 {
					m.focusIndex = 1
					m.artistInput.Blur()
					m.albumInput.Focus()
				} else {
					m.focusIndex = 0
					m.albumInput.Blur()
					m.artistInput.Focus()
				}
			case "enter":
				m.state = stateSearch
				m.error = nil
				return m, func() tea.Msg {
					res, err := m.curator.SearchMB(m.artistInput.Value(), m.albumInput.Value())
					if err != nil {
						return errorMsg(err)
					}
					return searchResultMsg(res)
				}
			}
			var cmd tea.Cmd
			if m.focusIndex == 0 {
				m.artistInput, cmd = m.artistInput.Update(msg)
			} else {
				m.albumInput, cmd = m.albumInput.Update(msg)
			}
			return m, cmd

		case stateResults:
			if msg.String() == "enter" {
				selected := m.list.SelectedItem().(item).release
				m.state = stateTagging
				m.status = fmt.Sprintf("Tagging files in %s...", m.targetDir)
				return m, func() tea.Msg {
					err := m.curator.TagDirectory(m.targetDir, selected.ID, nil)
					if err != nil {
						return errorMsg(err)
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

	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Bold(true).Render("PICKER: Album Tagger")

	switch m.state {
	case stateInput:
		s = fmt.Sprintf(
			"Target Directory: %s\n\n%s\n%s\n\n(tab to switch, enter to search)\n",
			m.targetDir,
			m.artistInput.View(),
			m.albumInput.View(),
		)
		if m.error != nil {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(fmt.Sprintf("\nError: %v\n", m.error))
		}

	case stateSearch:
		s = "\n  Searching MusicBrainz...\n"

	case stateResults:
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

func (i item) Title() string       { return i.release.Title }
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
