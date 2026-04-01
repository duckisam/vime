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

func HandleNormalInput(input string, m Model) (tea.Model, tea.Cmd) {
	m.commandOutput = ""
	var cmd tea.Cmd = nil
	switch input {
	case config.Quit, "ctrl+c":
		return m, tea.Quit
	case config.KeyDown, "down":
		m.cursor++
		SanitizeCursor(&m.cursor, len(m.entriesToDisplay))
		if m.cursor >= m.viewOffset + m.height - config.DisplayEntryOffset {
			m.viewOffset++
		}

	case config.KeyUp, "up":
		if m.cursor > 0{
			m.cursor = max(m.cursor - 1, 0)
			if m.cursor < m.viewOffset {
				m.viewOffset = max(m.viewOffset - 1, 0)
			}
			if m.viewOffset > m.cursor{
				m.viewOffset = m.cursor
			}

		}

	case config.Next:
		if m.entriesToDisplay[m.cursor].IsDir() {
			m.path = filepath.Join(m.path, m.entriesToDisplay[m.cursor].Name()) + "/"
			m.viewOffset = 0
			m.cursor = 0
			m.filter = ""
			cmd = loadDir(m.path)
		}

	case config.Back:
		if !(m.path == "/") {
			m.path = explorer.PathWalkBack(m.path)
			m.viewOffset = 0
			m.cursor = 0
			m.filter = ""
			cmd = loadDir(m.path)
		}

	case config.Confirm:
		m.filter = ""
		if m.entriesToDisplay[m.cursor].IsDir() {
			m.path = filepath.Join(m.path, m.entriesToDisplay[m.cursor].Name()) + "/"
			m.cursor = 0
			m.viewOffset = 0
			cmd = loadDir(m.path)
		} else {
			cmd := exec.Command(config.EditorCommand, m.path + m.entriesToDisplay[m.cursor].Name())
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return nil
			})
		}

	case config.CommandCopyPath:
		m.commandOutput, _ = ParseCommand("copy_string "+m.path+m.entriesToDisplay[m.cursor].Name(), m)

	case config.CommandModeToggle:
		m.mode = ModeCommand
		m.input.Focus()
		m.input.Prompt = config.CommandModeToggle

	case config.CommandCreateDir:
		m.mode = ModeCommand
		m.input.Focus()
		m.input.Prompt = config.CommandModeToggle
		m.input.SetValue("create_dir " + m.path)

	case config.CommandCreateFile:
		m.mode = ModeCommand
		m.input.Focus()
		m.input.Prompt = config.CommandModeToggle
		m.input.SetValue("create_file " + m.path)
	
	
	case config.CommandRemove:
		m.mode = ModeCommand
		m.input.Focus()
		m.input.Prompt = config.CommandModeToggle
		m.input.SetValue("remove " + m.path + m.entriesToDisplay[m.cursor].Name())

	case config.CommandRename:
		m.mode = ModeCommand
		m.input.Focus()
		m.input.Prompt = config.CommandModeToggle
		fileExt := filepath.Ext(m.entriesToDisplay[m.cursor].Name())
		command := "rename " + m.path + m.entriesToDisplay[m.cursor].Name() + " " + m.path + fileExt
		m.input.SetValue(command)
		m.input.SetCursor(len(command) - len(fileExt))

	case config.NormalSearch:
		m.mode = ModeSearch
		m.input.Focus()
		m.input.Prompt = config.NormalSearch

	case config.FuzzySearch:
		m.mode = ModeSearch
		m.fzf = true
		m.input.Focus()
		m.input.Prompt = "fzf/"
	case "esc", "\x1b":
		m.filter = ""
		cmd = loadDir(m.path)
	}

	return m, cmd
}

func HandleCommandInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case config.Confirm:
			m.commandOutput, cmd = ParseCommand(m.input.Value(), m)
			m.input.Reset()
			m.mode = ModeNormal
			m.input.Blur()
			return m, cmd
		case "esc", "\x1b":
			m.input.Reset()
			m.mode = ModeNormal
			m.input.Blur()
			return m, nil
		case "backspace":
			if m.input.Value() == ""{
				m.input.Reset()
				m.mode = ModeNormal
				m.input.Blur()
			}
		case "tab":
			completed := autoComplete(m.input.Value(), m.path)
			m.input.SetValue(completed)
			m.input.CursorEnd()
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func HandleSeachInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case config.Confirm:
			m.input.Reset()
			m.mode = ModeNormal
			m.input.Blur()
			return m, nil
		case "esc", "\x1b":
			m.input.Reset()
			m.mode = ModeNormal
			m.input.Blur()
		case "backspace":
			if m.input.Value() == ""{
				m.input.Reset()
				m.mode = ModeNormal
				m.input.Blur()
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	m.filter = m.input.Value()
	m.viewOffset = 0
	m.cursor = 0

	if m.fzf {
		m.entriesToDisplay = FuzzySearch(m.filter, m.entries)
	} else {
		m.entriesToDisplay = NormalSearch(m.filter, m.entries)
	}

	return m, cmd
}

func ParseCommand(command string, m Model) (string, tea.Cmd) {
	commandParts := strings.Split(command, " ")
	var err error
	var cmd tea.Cmd = nil
	switch commandParts[0] {
	case "mv", "move", "rename", "rn":
		if len(commandParts) != 3 {
			err = errors.New("invaild args amount")
			break
		}

		_, err = os.Stat(commandParts[2])
		if err == nil{
			err = errors.New("file \"" + commandParts[2] + " already exists")
			break
		}

		err = os.Rename(commandParts[1], commandParts[2])
		cmd = loadDir(m.path)

	case "cp", "copy":
		switch os := runtime.GOOS; os {
		case "linux":
		}
		
	case "cs", "copy_string":
		clip.WriteAll(commandParts[1])
		return config.SuccessStyle.Render("string \"" + commandParts[1] + "\" copied to clipboard"), nil

	case "q", "quit", "Q":
		cmd = tea.Quit

	case "edit":
			command := exec.Command(config.EditorCommand, m.path)
			cmd = tea.ExecProcess(command, func(err error) tea.Msg {
				return nil
			})

	case "cd", "change_dir":
		if len(commandParts) <= 1{
			err = errors.New("No path provided")
			break
		}

		commandParts[1] = strings.TrimSpace(commandParts[1])

		if !explorer.IsVaildOsPath(commandParts[1]){
			err = errors.New("invaild path \"" + commandParts[1]+ "\"")
			break
		}

		if commandParts[1] == ".." {
			m.path = explorer.PathWalkBack(m.path)
		} else {
			m.path = explorer.ExpandPath(commandParts[1])
		}

		if !strings.HasSuffix(m.path, "/") {
			m.path += "/"
		}

		cmd = loadDir(m.path)

	case "create_file":
		if len(commandParts) <= 1{
			err = errors.New("no file name given")
			break
		}else if len(commandParts) > 2{
			err = errors.New("to many args")
			break
		}

		_, err = os.Stat(commandParts[1])
		if err == nil{
			err = errors.New("file \"" + commandParts[1] + " already exists")
			break
		}
		
		var file *os.File
		file, err = os.Create(commandParts[1])
		if err == nil{
			file.Close()
			cmd = loadDir(m.path)
		}

	case "create_dir":
		if len(commandParts) <= 1{
			err = errors.New("no dir name given")
			break
		}else if len(commandParts) > 2{
			err = errors.New("to many args")
			break
		}

		err = os.Mkdir(commandParts[1], 0755)
		cmd = loadDir(m.path)
	
	case "rm", "remove":
		if len(commandParts) <= 1{
			err = errors.New("no file/dir name given")
			break
		}else if len(commandParts) > 2{
			err = errors.New("to many args")
			break
		}

		err = os.Remove(commandParts[1])
		cmd = loadDir(m.path)

	default:
		return config.ErrorStyle.Render("command: " + command + " is not a command"), cmd
	}
	if err != nil {
		return config.ErrorStyle.Render("command: " + command + " was not successfull with error: " + err.Error()), cmd
	}
	return config.SuccessStyle.Render("command: " + command + " was successfull"), cmd
}

func SanitizeCursor(pos *int, length int) {
	if length == 0{
		*pos = 0
		return
	}

	if *pos < 0 {
		*pos = 0
	}

	if *pos >= length {
		*pos = length - 1
	}

}
