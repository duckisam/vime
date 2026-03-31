package explorer

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func FormatDirEntries(entries []fs.DirEntry) []fs.DirEntry{
	var dirs []fs.DirEntry
	var files []fs.DirEntry

	for _, entry := range entries{
		if entry.IsDir(){
			dirs = append(dirs, entry)
		}else{
			files = append(files, entry)
		}

	}
	return append(dirs, files...)
	
}

func PathWalkBack(path string) string{
	path = strings.TrimSpace(path)

	if path[len(path)-1] == '/'{
		path = path[:len(path)-1]
	}

	index := strings.LastIndex(path, "/")
	path = path[:index + 1]
	
	return path 
}

func ExpandPath(path string) string{
	if strings.HasPrefix(path, "~"){
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}

	abs, err := filepath.Abs(path)
	if err != nil{
		return path
	}

	return abs
}

func IsVaildOsPath(path string) bool {
	_, err := os.Stat(ExpandPath(path))
	return err == nil
}
