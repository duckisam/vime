package ui

import (
	"io/fs"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/duckisam/vime/internal/config"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	explorer "github.com/duckisam/vime/internal/explorer"
)

var LastPath string

type Mode int

const (
	ModeNormal Mode = iota
	ModeCommand
	ModeSearch
)

type Model struct {
	path             string
	entries          []fs.DirEntry
	height           int
	width            int
	cursor           int
	viewOffset       int
	mode             Mode
	input            textinput.Model
	commandOutput    string
	fzf              bool
	filter           string
	entriesToDisplay []fs.DirEntry
}

type dirLoadMsg struct {
	entries []fs.DirEntry
	path    string
}

func loadDir(path string) tea.Cmd {
	return func() tea.Msg {
		err := os.Chdir(path)
		if err != nil {
			return nil
		}

		dir, _ := os.Getwd()
		entries, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		return dirLoadMsg{entries: explorer.FormatDirEntries(entries), path: path}
	}
}

func (m Model) Init() tea.Cmd {
	return loadDir(m.path)
}

func New(initPath string) Model {
	ti := textinput.New()
	ti.Prompt = ":"
	return Model{
		path:    initPath,
		entries: nil,
		filter:  "",
		fzf:     false,
		input:   ti,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dirLoadMsg:
		m.entries = msg.entries
		m.path = msg.path
		m.cursor = 0
		m.viewOffset = 0
		if m.fzf {
			m.entriesToDisplay = FuzzySearch(m.filter, m.entries)
		} else {
			m.entriesToDisplay = NormalSearch(m.filter, m.entries)
		}
		SanitizeCursor(&m.cursor, len(m.entriesToDisplay))

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		SanitizeCursor(&m.cursor, len(m.entriesToDisplay))
		
		switch m.mode {
		case ModeCommand:
			return HandleCommandInput(msg, m)
		case ModeSearch:
			return HandleSeachInput(msg, m)
		default:
			return HandleNormalInput(msg.String(), m)
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	
	SanitizeCursor(&m.cursor, len(m.entriesToDisplay))
	end := min(m.viewOffset+m.height - config.DisplayEntryOffset, len(m.entriesToDisplay))

	for i := m.viewOffset; i < end; i++ {
		if i == m.cursor {
			s.WriteString(config.SelectedStyle.Render("> ") + formatEntry(m, i) + "\n")
		} else {
			s.WriteString("  " + formatEntry(m, i) + "\n")
		}
	}
	
	if len(m.entriesToDisplay) == 0 && len(m.entries) != 0 {
		s.WriteString(config.ErrorStyle.Render("invaild search: \"" + m.filter + "\""))
	}

	if len(m.entries) == 0 {
		s.WriteString(config.ErrorStyle.Render("directory \"" + m.path + "\" " + " is empty"))
	}

	for i := end - m.viewOffset; i < m.height - config.DisplayEntryOffset; i++ {
		s.WriteString("\n")
	}
	
	return gloss.JoinVertical(gloss.Left, s.String(), m.renderStatusBar(), m.renderCommandBar())
}
