package ui

import (
	"io/fs"
	"os"
	"strings"
	"github.com/duckisam/vime/internal/config"
 	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	explorer "github.com/duckisam/vime/internal/explorer"
)

var LastPath string

type Mode int

const (
	ModeNormal Mode = iota
	ModeCommand
)

type Model struct{
	path string
	entries []fs.DirEntry
	height int
	width int
	cursor int
	mode Mode
	input textinput.Model
	commandOutput string
}

type dirLoadMsg struct{
	entries []fs.DirEntry
}

func loadDir(path string) tea.Cmd{
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
			if(err != nil){
			return nil
		}
		return dirLoadMsg{entries: explorer.FormatDirEntries(entries)}
	}
}

func (m Model) Init() tea.Cmd{
	return loadDir(m.path)
}

func New(initPath string) Model {
	ti := textinput.New()
	ti.Placeholder = "Type command..."
	ti.Prompt = ":"
	return Model{
		path: initPath,
		entries: nil,
		input: ti,
	}
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dirLoadMsg:
		m.entries = msg.entries
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.mode == ModeCommand {
			return HandleCommandInput(msg, m)
		}
		return HandleNormalInput(msg.String(), m)
    }
    return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	
	for i := 0; i < len(m.entries); i++{
		if i == m.cursor{
			s.WriteString(config.SelectedStyle.Render("> ") + formatEntry(m, i) + "\n")
		}else{
			s.WriteString("  " + formatEntry(m, i) + "\n")
		}
	}

	if len(m.entries) == 0{
		s.WriteString("empty dir")
	}

	for i := len(m.entries); i < m.height - 2; i++{
		s.WriteString("\n")
	}

	return gloss.JoinVertical(gloss.Left, s.String(), m.renderStatusBar())
}
