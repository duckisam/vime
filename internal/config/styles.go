package config

import (
	gloss "github.com/charmbracelet/lipgloss"
)

const (
	selectedColor = "#FFB845" 
	errorColor    = "#f0193e"
	successColor  = "#3bd723"
)


var (
	SelectedStyle = gloss.NewStyle().Foreground(gloss.Color(selectedColor))
	ErrorStyle    = gloss.NewStyle().Foreground(gloss.Color(errorColor))
	SuccessStyle  = gloss.NewStyle().Foreground(gloss.Color(successColor))
)


