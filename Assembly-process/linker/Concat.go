package linker

import (
	"fmt"
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

		data, err := os.ReadFile(handlePossibleStdlibFilepath(filePath))
		if err != nil {
			return nil, nil, err
		}
		dir := filepath.Dir(filePath)
		stringData := string(data)

		for _, line := range strings.Split(stringData, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, includeSignifier) {
				break
			}
			// line should contain the relative path from the line data location to the include data
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
	locations = make([]uint16, uniquePaths.Size()+2)
	filePaths = make([]string, uniquePaths.Size()+2)
	filePaths[0] = filePathSave
	locations[0] = 0

	for i, val := range uniquePaths.Get() {
		dir := filepath.Dir(val)
		fmt.Println("updating dir: ", dir)
		if dir == stdLibLocation {
			fmt.Println("is std")
			p := handlePossibleStdlibFilepath(val)
			fmt.Println(p)
			filePaths[i+1] = p
			locations[i+1] = ProgramStdLibStart
		} else {
			filePaths[i+1] = val
			locations[i+1] = 0
		}
	}

	IncludeBase(&filePaths, &locations)

	return filePaths, locations, nil
}

func IncludeBase(filePaths *[]string, locations *[]uint16) {
	locationsDe := *locations
	filepathsDe := *filePaths

	rootPath := helper.GetRootPath()
	tablePath := rootPath + includeLocationUse + "/interruptTable.obj"
	taskPath := rootPath + includeLocationUse + "/scheduling.obj"
	locationHeader := rootPath + includeLocationUse + "/sys_location.h"

	locationsDe = append(locationsDe, InterrupttableLoc, 0x2381, 0)
	filepathsDe = append(filepathsDe, tablePath, taskPath, locationHeader)

}

func IncludeHeaders() (filePaths []string, locations []uint16, e error) {
	filePaths = make([]string, 0)
	return filePaths, locations, nil
}

func handlePossibleStdlibFilepath(filename string) string {
	dir := filepath.Dir(filename)
	fmt.Println(filename)
	if dir != stdLibLocation {
		fmt.Println("not stdlib")
		return filename
	}

	root := helper.GetRootPath()
	file := filepath.Join(filepath.Join(root, stdLibLocationUse), filepath.Base(filename))
	return file
}
