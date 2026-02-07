package fileHandling

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AirCraft009/mcc"
	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

// FindIncludes
// finds all include statements inside the file
// only accepts includes at the start
// after a line that isn't an include statement it returns\
//
// returns all unique include filepaths
func FindIncludes(filePath string, fsHelp *mcc.FSHelper) (filePaths []string, locations []uint16, e error) {
	filePathSave := filePath
	uniquePaths := helper.NewSet[string]()
	nextPaths := helper.NewQueue[string]()

	for filePath != "" {

		data, err := fsHelp.ResolveReadFile(filePath)
		if err != nil {
			return nil, nil, err
		}
		dir := filepath.Dir(filePath)
		stringData := string(data)

		for _, line := range strings.Split(stringData, "\n") {
			if !strings.HasPrefix(line, helper.IncludeSignifier) {
				continue
			}

			line = strings.TrimSpace(line)
			// line should contain the relative path from the line data location to the include data
			line = strings.TrimSpace(strings.TrimPrefix(line, helper.IncludeSignifier))
			cleanedPath := filepath.Clean(filepath.Join(dir, line))
			if !uniquePaths.IsExist(cleanedPath) {
				nextPaths.Enqueue(cleanedPath)
				uniquePaths.Add(cleanedPath)
			}
		}
		// Dequeue == "" if nextPaths is empty
		filePath = nextPaths.Dequeue()
	}
	// add 1 to leave space for filePaths[0] = filePathsSave
	locations = make([]uint16, uniquePaths.Size()+1)
	filePaths = make([]string, uniquePaths.Size()+1)
	filePaths[0] = filePathSave

	for i, val := range uniquePaths.Get() {
		filePaths[i+1] = val
	}

	includeBaseComponents(&filePaths, &locations)

	return filePaths, locations, nil
}

// includeBaseComponents
//
// This includes everything that is always included.
//
// It places the interrupttable at the correct position.
// It includes the scheduler and all headers in lib/globalHeaders
func includeBaseComponents(filePaths *[]string, locations *[]uint16) {
	locationsDe := *locations
	filepathsDe := *filePaths

	rootPath := helper.GetRootPath()
	tablePath := filepath.Join(rootPath, filepath.Join(helper.IncludeLocationUse, "/interruptTable.obj"))
	taskPath := filepath.Join(rootPath, filepath.Join(helper.IncludeLocationUse, "/scheduling.obj"))
	IncludeHeaders(filePaths, locations)

	locationsDe = append(locationsDe, pkg.Interrupttable, 0)
	filepathsDe = append(filepathsDe, tablePath, taskPath)

	*locations = locationsDe
	*filePaths = filepathsDe
}

// IncludeHeaders
//
// adds all headers in lib/globalHeaders to the filepaths
func IncludeHeaders(filePaths *[]string, locations *[]uint16) {
	locationsDe := *locations
	filepathsDe := *filePaths

	rootPath := helper.GetRootPath()
	headerPath := rootPath + helper.GlobalHeaderLocation
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
