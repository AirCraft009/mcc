package linker

import (
	"mcc/helper"
	"os"
	"path/filepath"
	"strings"
)

// FindIncludes
// finds all include statements inside the file
// only accepts includes at the start
// after a line that isn't an include statement it returns\
//
// returns all unique include filepaths
func FindIncludes(filePath string) (filePaths []string, locations []uint16, e error) {
	filePathSave := filePath
	uniquePaths := helper.NewSet[string]()
	nextPaths := helper.NewQueue[string]()

	for filePath != "" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, err
		}
		dir := filepath.Dir(filePath)
		stringData := string(file)

		for _, line := range strings.Split(stringData, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, includeSignifier) {
				break
			}
			// line should contain the relative path from the line file location to the include file
			line = strings.TrimSpace(strings.TrimPrefix(line, includeSignifier))
			cleanedPath := filepath.Clean(filepath.Join(dir, line))
			if !uniquePaths.IsExist(cleanedPath) {
				nextPaths.Enqueue(cleanedPath)
				uniquePaths.Add(cleanedPath)
			}
		}
		// "" if empty
		filePath = nextPaths.Dequeue()
	}
	locations = make([]uint16, uniquePaths.Size()+1)
	filePaths = make([]string, uniquePaths.Size()+1)
	filePaths[0] = filePathSave
	locations[0] = 0
	for i, val := range uniquePaths.Get() {
		dir := filepath.Dir(val)
		if dir == stdLibLocation {
			filePaths[i+1] = filepath.Join(filepath.Clean(stdLibLocation), filepath.Base(val))
			locations[i+1] = ProgramStdLibStart
		} else {
			filePaths[i+1] = val
			locations[i+1] = 0
		}
	}

	return filePaths, locations, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
