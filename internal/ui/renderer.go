package ui

import (
	"strings"
	gloss "github.com/charmbracelet/lipgloss"
	config "github.com/duckisam/vime/internal/config"
	icons "github.com/epilande/go-devicons"
)

func (m *Model) renderStatusBar() string{
	barStyle := gloss.NewStyle().
	Background(gloss.Color("##00000000")).
	Width(m.width).
	Padding(0, 1)

	left := m.path

	if left == "/"{
		left = ""
	}else if m.mode == ModeCommand{
		left = m.input.View() 
	}else if m.commandOutput != ""{
		left = m.commandOutput
	}

	right := config.Back + " back  " + config.Confirm + " open  " + config.Quit + " quit"

	gap := max(m.width - gloss.Width(left) - gloss.Width(right) - 2, 0)

	bar := left + strings.Repeat(" ", gap) + right
	return barStyle.Render(bar)
}

func formatEntry(m Model, entryIndex int) string{
	icon := icons.IconForPath(m.path + m.entries[entryIndex].Name())
	iconStyle := gloss.NewStyle().Foreground(gloss.Color(icon.Color))
	
	var s strings.Builder
	entry := m.entries[entryIndex]
	entryString := entry.Name()

	if entry.IsDir(){
		entryString += "/"
	}

	s.WriteString(iconStyle.Render(icon.Icon + " "))
	if entryIndex == m.cursor {
		s.WriteString(config.SelectedStyle.Render(entryString))
	} else {
		s.WriteString(entryString)
	}
	
	return s.String()
}

