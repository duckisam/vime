package main

import (
	"fmt"
	os "os"
	tea "github.com/charmbracelet/bubbletea"
	ui "github.com/duckisam/vime/internal/ui"
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

	for _, cmd := range ui.QuitComands{
		cmd.Run()
	}

	fmt.Println(ui.LastPath)
}

