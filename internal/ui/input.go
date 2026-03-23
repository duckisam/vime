package ui

import (
	"errors"
	"os"
	"runtime"
	"os/exec"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	config "github.com/duckisam/vime/internal/config"
	explorer "github.com/duckisam/vime/internal/explorer"
)


func HandleNormalInput(input string, m Model) (tea.Model, tea.Cmd){
	m.commandOutput = ""
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
		case config.CommandModeToggle:
			m.mode = ModeCommand
			m.input.Focus()
			return m, nil

		}
		return m, nil
}

func HandleCommandInput(msg tea.Msg,  m Model) (tea.Model, tea.Cmd){
	var cmd tea.Cmd
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.String(){
		case config.Confirm:
			m.commandOutput, cmd = ParseCommand(m.input.Value(), m)
			m.input.Reset()
			m.mode = ModeNormal
			m.input.Blur()
			return m, cmd
		case "esc":
			m.input.Reset()
			m.mode = ModeNormal
			m.input.Blur()
			return m, nil
		}

	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func ParseCommand(command string, m Model) (string, tea.Cmd){
	commandParts := strings.Split(command, " ")
	var err error
	var cmd tea.Cmd = nil
	switch commandParts[0]{
	case "mv", "move", "rename", "rn":
		if len(commandParts) == 3{
			err = os.Rename(commandParts[1], commandParts[2])
			cmd = loadDir(m.path)
			break
		}
		err = errors.New("invaild args amount")
	case "cp", "copy":
		var copyComand exec.Cmd
		switch os := runtime.GOOS; os{
		case "linux":

			

		}


		
	default:
		return config.ErrorStyle.Render("command: " + command + " is not a command"), cmd
	}
	if err != nil{
		return  config.ErrorStyle.Render("command: " + command + " was not successfull with error: " + err.Error()), cmd
	}

	return config.SuccessStyle.Render("command: " + command + " was successfull"), cmd




}
