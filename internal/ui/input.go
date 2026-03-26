package ui

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	clip "github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	config "github.com/duckisam/vime/internal/config"
	explorer "github.com/duckisam/vime/internal/explorer"
)


func HandleNormalInput(input string, m Model) (tea.Model, tea.Cmd){
	m.commandOutput = ""
	var cmd tea.Cmd = nil
	switch input {
		case config.Quit, "ctrl+c":
			return m, tea.Quit
		case config.KeyDown, "down":
			m.cursor = min(m.cursor + 1, len(m.entries) - 1)

		case config.KeyUp, "up":
			m.cursor = max(m.cursor - 1, 0)
				
		case config.Next:
			if m.entries[m.cursor].IsDir(){
				m.path = filepath.Join(m.path, m.entries[m.cursor].Name()) + "/"
				m.viewOffset = 0
				m.cursor = 0
				cmd = loadDir(m.path)
			}
		case config.Back:
			if !(m.path == "/"){
				m.path = explorer.PathWalkBack(m.path)
				m.viewOffset = 0
				m.cursor = 0
				cmd = loadDir(m.path)
			}
		case config.Confirm:
			if m.entries[m.cursor].IsDir(){
				m.path = filepath.Join(m.path, m.entries[m.cursor].Name()) + "/"
				m.cursor = 0
				m.viewOffset = 0
				cmd = loadDir(m.path)

			}else{
				cmd := exec.Command(config.EditorCommand, m. path + m.entries[m.cursor].Name())
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return nil
				})
			}
		case config.CommandCopyPath:
			m.commandOutput, _ = ParseCommand("copy_string " + m.path + m.entries[m.cursor].Name(), m)

		case config.CommandModeToggle:
			m.mode = ModeCommand
			m.input.Focus()

		case config.NormalSearch:
			m.mode = ModeSearch
			m.input.Focus()
		}
		return m, cmd
}

func HandleCommandInput(msg tea.Msg,  m Model) (tea.Model, tea.Cmd){
	var cmd tea.Cmd
	m.input.Prompt = config.CommandModeToggle
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

func HandleSeachInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd){
	var cmd tea.Cmd
	m.input.Prompt = "/"
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.String(){
		case config.Confirm:
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
	m.filter = m.input.Value()

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
		switch os := runtime.GOOS; os{
		case "linux":
		}
	case "cs", "copy_string":
		clip.WriteAll(commandParts[1])
		return config.SuccessStyle.Render("string \"" + commandParts[1] + "\" copied to clipboard"), nil
		
	case "q", "quit", "Q":
		cmd = tea.Quit
	case "cd":
		if explorer.IsVaildOsPath(commandParts[1]){
			err = errors.New("invaild path")
			break
		}

		commandParts[1] = strings.TrimSpace(commandParts[1])

		if commandParts[1] == ".."{
			m.path = explorer.PathWalkBack(m.path)
		}else{
			m.path = commandParts[1]
		}

		m.path = explorer.ExpandPath(commandParts[1])
		if !strings.HasSuffix(m.path, "/"){
			m.path += "/"
		}
		
		cmd = loadDir(m.path)
		
	default:
		return config.ErrorStyle.Render("command: " + command + " is not a command"), cmd
	}
	if err != nil{
		return  config.ErrorStyle.Render("command: " + command + " was not successfull with error: " + err.Error()), cmd
	}
	return config.SuccessStyle.Render("command: " + command + " was successfull"), cmd
}

func SanitizeCursor(pos int, length int) int{
	if pos < 0{
		return 0
	}

	if pos >= length{
		return  length -1
	}

	return pos
	
}
