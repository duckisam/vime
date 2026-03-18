package ui

import (
	"strings"
	"os"
	fs "io/fs"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	config "github.com/duckisam/vime/internal/config"
	explorer "github.com/duckisam/vime/internal/explorer"
)

type Model struct{
	path string
	entries []fs.DirEntry
	height int
	width int
	cursor int
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
	return Model{
		path: initPath,
		entries: nil,
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
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case config.KeyDown, "down":
			if m.cursor < len(m.entries) - 1{
				m.cursor++
			}
		case config.KeyUp, "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case config.Confirm:
			if m.entries[m.cursor].IsDir(){
				m.path +=  m.entries[m.cursor].Name() + "/"
				return m, loadDir(m.path)
			}
		case config.Back:
			m.path = explorer.PathWalkBack(m.path)
			return m, loadDir(m.path)
		}
    }
    return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	for i, entry := range m.entries{
		if i == m.cursor{
			s.WriteString("> " + fs.FormatDirEntry(entry) + "\n")
		}else{
			s.WriteString("  " + fs.FormatDirEntry(entry) + "\n")
		}
	}

	for i := len(m.entries); i < m.height - 2; i++{
		s.WriteString("\n")
	}

	return gloss.JoinVertical(gloss.Left, s.String(), m.renderStatusBar())
}
