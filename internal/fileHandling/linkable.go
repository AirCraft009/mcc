package fileHandling

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/AirCraft009/mcc"
	loader "github.com/AirCraft009/mcc/internal/asm/assembly"
	"github.com/AirCraft009/mcc/internal/asm/assembly/assembler"
	"github.com/AirCraft009/mcc/pkg"

	"golang.org/x/sync/errgroup"
)

const (
	ObjectF = iota + 1
	AsmF
	HeaderF
	CSOurceF
)

var acceptedFiletypes = map[string]byte{
	".asm": AsmF,
	".obj": ObjectF,
	".h":   HeaderF,
	".c":   CSOurceF,
}

type Linkables struct {
	Files []*LinkFile
	size  uint32
	write bool
	mutex sync.Mutex
}

type LinkFile struct {
	Path     string
	Data     []byte
	Location uint16
	FileT    uint8
}

func NewLinkables(size int) *Linkables {
	return &Linkables{make([]*LinkFile, size), 0, true, sync.Mutex{}}
}

// AddArraysMultiThreaded
// sadly unusable - as it causes nondeterministic output by using goroutines
func (link *Linkables) AddArraysMultiThreaded(filePaths []string, locations []uint16, fsHelper *mcc.FSHelper) error {
	//clears out anything should be empty anyway but who knows
	link.Files = make([]*LinkFile, len(filePaths))

	for _, filePath := range filePaths {
		if filePath == "" {
			fmt.Println("Empty file path")
		}
		fmt.Printf("Adding %s file\n", filePath)
	}

	if len(filePaths) != len(locations) {
		return errors.New("AddArrays did not receive the same number of files and locations")
	}

	g := &errgroup.Group{}
	//fmt.Println(len(filePaths))
	//fmt.Println(filePaths)
	for i := 0; i < len(filePaths); i++ {
		file, location := filePaths[i], locations[i]

		f, l, j := file, location, i
		//fmt.Println("sent gothread:", i, file)
		g.Go(func() error {
			err := link.addFileMultiThreaded(f, l, j, fsHelper)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return g.Wait()
}

func (link *Linkables) addFileMultiThreaded(filePath string, location uint16, index int, fsHelper *mcc.FSHelper) (err error) {
	ext := filepath.Ext(filePath)
	fileT := acceptedFiletypes[ext]
	if fileT == 0 {
		//fmt.Println("problem file: ", filePath, index)
		log.Printf("Unsupported file type: %s treating as ASM source file\n", ext)
		fileT = AsmF
	}

	data, err := fsHelper.ResolveReadFile(filePath)
	if err != nil {
		return err
	}

	linkF := &LinkFile{Path: filePath, Data: data, FileT: acceptedFiletypes[ext], Location: location}
	link.Files[index] = linkF
	atomic.AddUint32(&link.size, 1)

	return nil
}

func (link *Linkables) GetFiles() []*LinkFile {
	files := make([]*LinkFile, len(link.Files))
	copy(files, link.Files[:link.size])
	fmt.Println("fileOutput: ", files, len(link.Files))
	return files
}

func (link *Linkables) formatObjectFiles(outPath string, write bool, logger *log.Logger) (objectFiles map[*pkg.ObjectFile]uint16, err error) {
	locations := make(map[uint16]uint16)
	objFiles := make(map[*pkg.ObjectFile]uint16)

	for _, file := range link.Files {

		logger.Printf("Handling file %s\n", file.Path)

		if file.FileT == HeaderF {
			continue
		}

		var objFile *pkg.ObjectFile
		if file.FileT == ObjectF {
			objFile, err = pkg.FormatObjectFile(file.Data)
			if err != nil {
				return nil, err
			}

		} else {
			objFile = assembler.AssembleAndWrite(string(file.Data), outPath, write, logger)
		}

		// is location already used by another file
		if value, ok := locations[file.Location]; ok {
			logger.Printf("Location conflict %s\n", file.Path)

			objFiles[objFile] = file.Location + value
			locations[file.Location] = uint16(len(objFile.Code)) + value

		} else {
			objFiles[objFile] = file.Location
			locations[file.Location] = uint16(len(objFile.Code))
		}

	}

	return objFiles, nil
}

func (link *Linkables) GetObjectFiles(outPath string, write bool, logger *log.Logger) (objectFiles map[*pkg.ObjectFile]uint16, data []byte, err error) {
	objectFiles, err = link.formatObjectFiles(outPath, write, logger)
	if err != nil {
		return nil, nil, err
	}

	var objF = make([]*pkg.ObjectFile, 0, len(objectFiles))
	for obj, _ := range objectFiles {
		objF = append(objF, obj)
	}

	data = loader.LoadData(objF)

	return objectFiles, data, nil
}
