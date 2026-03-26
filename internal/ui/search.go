package ui

import (
	"io/fs"
	"strings"
)



func NormalSearch(filter string, targets []fs.DirEntry) []fs.DirEntry{
	if filter == ""{
		return targets
	}

	var results []fs.DirEntry
	for _, target := range targets{
		if strings.Contains(target.Name(), filter){
			results = append(results, target)
		}
	}

	return results


}


func FuzzySearch(pattern string, targets []fs.DirEntry) []fs.DirEntry{
	if pattern == ""{
		return targets
	}


	var results []fs.DirEntry
	for _, target := range targets{
		patternIdx := 0
		for i := 0; i < len(target.Name()) && patternIdx < len(pattern); i++{
			if target.Name()[i] == pattern[patternIdx]{
				patternIdx++
			}
		}
		
		if patternIdx == len(pattern){
			results = append(results, target)
		}

	}
	return results
}

