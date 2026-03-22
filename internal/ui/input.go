package ui

import (
	"os/exec"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	config "github.com/duckisam/vime/internal/config"
	explorer "github.com/duckisam/vime/internal/explorer"
)

func HandleInput(input string, m Model) (tea.Model, tea.Cmd){
		switch strings.ToLower(input) {
		case config.Quit, "ctrl+c":
			LastPath = m.path
			return m, tea.Quit
		case config.KeyDown, "down":
			if m.cursor < len(m.entries) - 1{
				m.cursor++
			}
		case config.KeyUp, "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case config.Next:
			if m.entries[m.cursor].IsDir(){
				m.path +=  m.entries[m.cursor].Name() + "/"
				m.cursor = 0
				return m, loadDir(m.path)
			}
		case config.Back:
			if !(m.path == "/"){
				m.path = explorer.PathWalkBack(m.path)
				m.cursor = 0
				return m, loadDir(m.path)
			}
		case config.Confirm:
			if m.entries[m.cursor].IsDir(){
				m.path +=  m.entries[m.cursor].Name() + "/"
				m.cursor = 0
				return m, loadDir(m.path)
			}else{
				cmd := exec.Command(config.EditorCommand, m. path + m.entries[m.cursor].Name())
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return nil
				})
			}
		case config.ChangeDir:
			QuitComands = append(QuitComands, *exec.Command("bash", "cd"))
			return m, tea.Quit

		}
		return m, nil

}
