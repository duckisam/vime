package ui

import (
	"strings"
	gloss "github.com/charmbracelet/lipgloss"
	config "github.com/duckisam/vime/internal/config"
)

func (m Model) renderStatusBar() string{
	barStyle := gloss.NewStyle().
	Background(gloss.Color("##00000000")).
	Foreground(gloss.Color("#FAFAFA")).
	Width(m.width).
	Padding(0, 1)

	left := m.path
	right := config.Back + " back  " + config.Confirm + " open  " + config.Quit + " quit"

	gap := max(m.width - gloss.Width(left) - gloss.Width(right) - 2, 0)

	bar := left + strings.Repeat(" ", gap) + right
	return barStyle.Render(bar)
}

