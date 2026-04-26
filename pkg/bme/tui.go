package bme

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Messages
type StatusMsg struct {
	Component string
	Device    string
	Status    string
}

type ProgressMsg struct {
	Component string
	Device    string
	Percent   float64
}

type MBMatchesMsg struct {
	MBID     string
	Releases []mb_release
}

type MBSelectedMsg struct {
	Release mb_release
}

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(44)

	ripperStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	encoderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	taggerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	systemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))

	logStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1)
)

type model struct {
	ripperStatuses   map[string]string
	ripperProgresses map[string]progress.Model
	encoderStatus    string
	taggerStatus     string
	paranoiaMode     string

	encoderProgress progress.Model
	taggerProgress  progress.Model

	mbMatches []mb_release
	mbList    list.Model
	selecting bool

	logs   []string
	width  int
	height int
}

func NewModel() model {
	return model{
		ripperStatuses:   make(map[string]string),
		ripperProgresses: make(map[string]progress.Model),
		encoderStatus:    "Idle",
		taggerStatus:     "Idle",
		paranoiaMode:     GetParanoiaName(),
		encoderProgress: progress.New(progress.WithDefaultGradient()),
		taggerProgress:  progress.New(progress.WithDefaultGradient()),
		logs:             make([]string, 0),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for k, v := range m.ripperProgresses {
			v.Width = 40
			m.ripperProgresses[k] = v
		}
		m.encoderProgress.Width = 40
		m.taggerProgress.Width = 40

	case StatusMsg:
		var styledStatus string
		var statusChanged bool

		switch msg.Component {
		case "ripper":
			if m.ripperStatuses[msg.Device] != msg.Status {
				statusChanged = true
			}
			m.ripperStatuses[msg.Device] = msg.Status
			if _, ok := m.ripperProgresses[msg.Device]; !ok {
				p := progress.New(progress.WithDefaultGradient())
				p.Width = 40
				m.ripperProgresses[msg.Device] = p
			}
			styledStatus = ripperStyle.Render(fmt.Sprintf("[%s:%s] %s", msg.Component, msg.Device, msg.Status))
		case "encoder":
			if m.encoderStatus != msg.Status {
				statusChanged = true
			}
			m.encoderStatus = msg.Status
			styledStatus = encoderStyle.Render(fmt.Sprintf("[%s] %s", msg.Component, msg.Status))
		case "tagger":
			if m.taggerStatus != msg.Status {
				statusChanged = true
			}
			m.taggerStatus = msg.Status
			styledStatus = taggerStyle.Render(fmt.Sprintf("[%s] %s", msg.Component, msg.Status))
		default:
			statusChanged = true // Always log system messages
			styledStatus = systemStyle.Render(fmt.Sprintf("[%s] %s", msg.Component, msg.Status))
		}

		if statusChanged && msg.Status != "Idle" {
			m.logs = append(m.logs, styledStatus)
			if len(m.logs) > 50 {
				m.logs = m.logs[1:]
			}
		}

	case ProgressMsg:
		switch msg.Component {
		case "ripper":
			if p, ok := m.ripperProgresses[msg.Device]; ok {
				cmds = append(cmds, p.SetPercent(msg.Percent))
				m.ripperProgresses[msg.Device] = p
			} else {
				p := progress.New(progress.WithDefaultGradient())
				p.Width = 40
				cmds = append(cmds, p.SetPercent(msg.Percent))
				m.ripperProgresses[msg.Device] = p
			}
		case "encoder":
			cmds = append(cmds, m.encoderProgress.SetPercent(msg.Percent))
		case "tagger":
			cmds = append(cmds, m.taggerProgress.SetPercent(msg.Percent))
		}

	case progress.FrameMsg:
		for k, v := range m.ripperProgresses {
			newRipper, ripperCmd := v.Update(msg)
			m.ripperProgresses[k] = newRipper.(progress.Model)
			if ripperCmd != nil {
				cmds = append(cmds, ripperCmd)
			}
		}

		newEncoder, encoderCmd := m.encoderProgress.Update(msg)
		m.encoderProgress = newEncoder.(progress.Model)
		if encoderCmd != nil {
			cmds = append(cmds, encoderCmd)
		}

		newTagger, taggerCmd := m.taggerProgress.Update(msg)
		m.taggerProgress = newTagger.(progress.Model)
		if taggerCmd != nil {
			cmds = append(cmds, taggerCmd)
		}

	case MBMatchesMsg:
		m.mbMatches = msg.Releases
		m.selecting = true
		items := make([]list.Item, len(msg.Releases))
		for i, r := range msg.Releases {
			items[i] = mbItem{r}
		}
		m.mbList = list.New(items, list.NewDefaultDelegate(), m.width-10, m.height-15)
		m.mbList.Title = fmt.Sprintf("Select Release for DiscID: %s", msg.MBID)

	case tea.KeyMsg:
		// Global keys that work in any state
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+l":
			return m, tea.ClearScreen
		}

		if m.selecting {
			switch msg.String() {
			case "enter":
				selected := m.mbList.SelectedItem().(mbItem).release
				m.selecting = false
				// Send back to tagger
				go func() {
					selectionChan <- selected
				}()
				return m, nil
			case "q":
				return m, tea.Quit
			}
			m.mbList, cmd = m.mbList.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "p":
			mode := ToggleParanoia()
			m.paranoiaMode = GetParanoiaName()
			m.logs = append(m.logs, systemStyle.Render(fmt.Sprintf("[system] Paranoia set to: %s", mode)))
		case "X":
			PurgeDirectories()
			m.logs = append(m.logs, systemStyle.Render("[system] Working directories purged"))
		case "q":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	header := headerStyle.Render("BME: Batch Music Encoder Dashboard")

	// Collect ripper views
	var ripperViews []string
	keys := make([]string, 0, len(m.ripperStatuses))
	for k := range m.ripperStatuses {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		status := m.ripperStatuses[k]
		prog := m.ripperProgresses[k]
		ripperViews = append(ripperViews, ripperStyle.Render(fmt.Sprintf("Ripper(%s): ", k))+status)
		ripperViews = append(ripperViews, prog.View())
	}
	if len(ripperViews) == 0 {
		ripperViews = append(ripperViews, ripperStyle.Render("Ripper:   ")+"Idle")
		ripperViews = append(ripperViews, m.encoderProgress.View()) // dummy to keep spacing? No, maybe just empty.
	}

	encoderView := encoderStyle.Render("Encoder:  ") + m.encoderStatus
	taggerView := taggerStyle.Render("Tagger:   ") + m.taggerStatus
	paranoiaView := systemStyle.Render("Paranoia: ") + m.paranoiaMode

	// Calculate heights
	mainHeight := m.height - 6
	if mainHeight < 10 {
		mainHeight = 10
	}
	// with multiple rippers, we might need more space, but let's see.

	statusElements := []string{}
	statusElements = append(statusElements, ripperViews...)
	statusElements = append(statusElements,
		encoderView,
		m.encoderProgress.View(),
		taggerView,
		m.taggerProgress.View(),
		"\n"+paranoiaView,
	)

	statusCol := statusStyle.
		Height(mainHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left, statusElements...))

	// Take last N lines of logs that fit in height
	availableLogLines := mainHeight - 2
	if availableLogLines < 1 {
		availableLogLines = 1
	}

	displayLogs := m.logs
	if len(displayLogs) > availableLogLines {
		displayLogs = displayLogs[len(displayLogs)-availableLogLines:]
	}
	logContent := strings.Join(displayLogs, "\n")

	logView := logStyle.
		Width(m.width - 49).
		Height(mainHeight).
		Render(logContent)

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, statusCol, logView)

	if m.selecting {
		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			"\n",
			m.mbList.View(),
		)
	}

	help := "\n[q] quit | [p] toggle paranoia | [ctrl+l] redraw | [X] purge work dirs"

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"\n",
		mainView,
		help,
	)
}

type mbItem struct {
	release mb_release
}

func (i mbItem) Title() string { return i.release.Title }
func (i mbItem) Description() string {
	return fmt.Sprintf("Tracks: %d | Artist: %s | Title: %s | Country: %s | Barcode: %s | Disambig: %s", len(i.release.Tracks), i.release.AlbumArtist, i.release.Title, i.release.Country, i.release.Barcode, i.release.Disambiguation)
}
func (i mbItem) FilterValue() string { return i.release.Title }
