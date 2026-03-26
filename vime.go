package main

import (
	"os"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	ui "github.com/duckisam/vime/internal/ui"
)

const version = "1.0.0"

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
	
	lastdir, _ := os.Getwd()

	fmt.Println(lastdir)
}

