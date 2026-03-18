package main

import (
	tea "github.com/charmbracelet/bubbletea"
	ui "github.com/duckisam/vime/internal/ui"
	os "os"
)

func main(){
	initPath, err := os.Getwd()

	if err != nil{
		panic(err)
	}

	initPath += "/"

	p := tea.NewProgram(ui.New(initPath), tea.WithAltScreen())

	if _, err := p.Run(); err != nil{
		panic(err)
	}
}

