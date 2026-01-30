package linker

import (
	"fmt"
	helper2 "mcc/internal/helper"
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
	uniquePaths := helper2.NewSet[string]()
	nextPaths := helper2.NewQueue[string]()

	for filePath != "" {

		data, err := os.ReadFile(handlePossibleStdlibFilepath(filePath))
		if err != nil {
			return nil, nil, err
		}
		dir := filepath.Dir(filePath)
		stringData := string(data)

		for _, line := range strings.Split(stringData, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, helper2.IncludeSignifier) {
				break
			}
			// line should contain the relative path from the line data location to the include data
			line = strings.TrimSpace(strings.TrimPrefix(line, helper2.IncludeSignifier))
			cleanedPath := filepath.Clean(filepath.Join(dir, line))
			if !uniquePaths.IsExist(cleanedPath) {
				nextPaths.Enqueue(cleanedPath)
				uniquePaths.Add(cleanedPath)
			}
		}
		// "" if empty
		filePath = nextPaths.Dequeue()
	}
	// add 1 to leave space for filePaths[0] = filePathsSave
	locations = make([]uint16, uniquePaths.Size()+1)
	filePaths = make([]string, uniquePaths.Size()+1)
	filePaths[0] = filePathSave
	locations[0] = 0

	for i, val := range uniquePaths.Get() {
		dir := filepath.Dir(val)
		//fmt.Println("updating dir: ", dir)
		if dir == helper2.StdLibLocation {
			//fmt.Println("is std")
			p := handlePossibleStdlibFilepath(val)
			fmt.Println(p)
			filePaths[i+1] = p
			locations[i+1] = helper2.ProgramStdLibStart
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

	rootPath := helper2.GetRootPath()
	tablePath := rootPath + helper2.IncludeLocationUse + "/interruptTable.obj"
	taskPath := rootPath + helper2.IncludeLocationUse + "/scheduling.obj"
	IncludeHeaders(filePaths, locations)

	locationsDe = append(locationsDe, helper2.InterrupttableLoc, 0x2381)
	filepathsDe = append(filepathsDe, tablePath, taskPath)

}

func IncludeHeaders(filePaths *[]string, locations *[]uint16) {
	locationsDe := *locations
	filepathsDe := *filePaths

	rootPath := helper2.GetRootPath()
	headerPath := rootPath + helper2.GlobalHeaderLocation
	dir, err := os.ReadDir(headerPath)
	if err != nil {
		panic(err.Error())
	}
	dirLen := len(dir)
	prevlen := len(locationsDe)
	locationsTemp := make([]uint16, dirLen+prevlen)
	filePathsTemp := make([]string, dirLen+prevlen)
	copy(locationsTemp, locationsDe)
	copy(filePathsTemp, filepathsDe)

	for i, file := range dir {
		if file.IsDir() {
			panic("Directory " + file.Name() + " is not a header and shouldn't be in the globalHeaders directory")
		}
		filePathsTemp[i+prevlen] = filepath.Join(headerPath, file.Name())
	}

	*locations = locationsTemp
	*filePaths = filePathsTemp
}

func handlePossibleStdlibFilepath(filename string) string {
	dir := filepath.Dir(filename)
	fmt.Println(filename)
	if dir != helper2.StdLibLocation {
		//fmt.Println("not stdlib")
		return filename
	}

	root := helper2.GetRootPath()
	file := filepath.Join(filepath.Join(root, helper2.StdLibLocationUse), filepath.Base(filename))
	return file
}
