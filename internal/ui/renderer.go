package ui

import (
	"os/user"
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
	config "github.com/duckisam/vime/internal/config"
	icons "github.com/epilande/go-devicons"
)

func (m Model) renderStatusBar() string{
	barStyle := gloss.NewStyle().
	Background(gloss.Color("#00000000")).
	Width(m.width)

	left := m.path
	usr, _ := user.Current()
	
	if strings.HasPrefix(left, usr.HomeDir){
		left = strings.Replace(left, usr.HomeDir, "~", 1)
	}

	right := config.Back + " back  " + config.Confirm + " open  " + config.Quit + " quit"

	gap := max(m.width - gloss.Width(left) - gloss.Width(right) - 2, 0)

	bar := left + strings.Repeat(" ", gap) + right
	return barStyle.Render(bar)
}

func (m Model) renderCommandBar() string{
	var toDisplay string
	switch m.mode{
	case ModeNormal:
		toDisplay = ""
	default:
		toDisplay = m.input.View()
	}

	if m.commandOutput != ""{
		toDisplay = m.commandOutput
	}
	

	return toDisplay
}

func formatEntry(m Model, entryIndex int) string{
	icon := icons.IconForPath(m.path + m.entriesToDisplay[entryIndex].Name())
	iconStyle := gloss.NewStyle().Foreground(gloss.Color(icon.Color))
	
	var s strings.Builder
	entry := m.entriesToDisplay[entryIndex]
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

