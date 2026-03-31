package ui

import (
	"os"
	"path/filepath"
	"strings"
	explorer "github.com/duckisam/vime/internal/explorer"
)

var commands = []string{
	"move", "rename", "copy_string", "quit", "edit",
	"remove", "create_file", "create_dir", "change_dir",
}

func autoComplete(input string, currentPath string) string {
	if input == "" {
		return ""
	}

	inputParts := strings.Split(input, " ")

	if len(inputParts) == 1 {
		return getCommands(inputParts[0])
	}

	lastArg := inputParts[len(inputParts)-1]
	completion := getPath(lastArg)

	return input + completion
}

func getCommands(input string) string {
	var filtered []string
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, input) {
			filtered = append(filtered, cmd)
		}
	}

	if len(filtered) == 0 {
		return input
	}
	if len(filtered) == 1 {
		return filtered[0]
	}

	common := filtered[0]
	for _, s := range filtered[1:] {
		common = longestCommonPrefix(common, s)
	}
	return common
}

func getPath(currentInput string) string {
	if currentInput == "" {
		return ""
	}

	expanded := explorer.ExpandPath(currentInput)
	
	dir := filepath.Dir(expanded)
	prefix := filepath.Base(expanded)

	if strings.HasSuffix(currentInput, "/") || currentInput == "." {
		dir = expanded
		prefix = ""
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var matches []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			if entry.IsDir() {
				name += "/"
			}
			matches = append(matches, name)
		}
	}

	if len(matches) == 0 {
		return ""
	}

	if len(matches) == 1 {
		return matches[0][len(prefix):]
	}

	common := matches[0]
	for _, m := range matches[1:] {
		common = longestCommonPrefix(common, m)
	}

	if len(common) > len(prefix) {
		return common[len(prefix):]
	}

	return ""
}

func longestCommonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[:i]
}
