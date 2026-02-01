package linker

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/AirCraft009/mcc/internal/assembler"
	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"

	"golang.org/x/sync/errgroup"
)

const (
	ObjectF = iota + 1
	AsmF
	HeaderF
)

var acceptedFiletypes = map[string]byte{
	".asm": AsmF,
	".obj": ObjectF,
	".h":   HeaderF,
}

type Linkables struct {
	Files []*LinkFile
	ptr   uint32
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
	return &Linkables{make([]*LinkFile, size), 0, 0, true, sync.Mutex{}}
}

// AddArraysMultiThreaded
// sadly unusable - as it causes nondeterministic output by using goroutines
func (link *Linkables) AddArraysMultiThreaded(filePaths []string, locations []uint16) error {
	//clears out anything should be empty anyway but who knows
	link.Files = make([]*LinkFile, len(filePaths))
	//force it to start adding at 0
	link.ResetPtr()

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
		//fmt.Println("adding file: ", f, i)
		g.Go(func() error {
			err := link.addFileMultiThreaded(f, l, j)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return g.Wait()
}

func (link *Linkables) addFileMultiThreaded(filePath string, location uint16, index int) (err error) {
	ext := filepath.Ext(filePath)

	if acceptedFiletypes[ext] == 0 {
		//fmt.Println("problem file: ", filePath, index)
		return errors.New("Unsupported file type: " + ext)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	linkF := &LinkFile{Path: filePath, Data: data, FileT: acceptedFiletypes[ext], Location: location}
	link.mutex.Lock()
	link.Files[index] = linkF
	link.ptr++
	link.size++
	link.mutex.Unlock()

	return nil
}

// GetFile
//
// file may be nil
func (link *Linkables) GetFile() (file *LinkFile, nonNil bool) {
	link.mutex.Lock()
	if link.ptr == 0 {
		link.mutex.Unlock()
		return nil, false
	}
	link.ptr = helper.SatSubU32(link.ptr, 1)
	file = link.Files[link.ptr]
	link.mutex.Unlock()
	return file, file != nil
}

func (link *Linkables) ResetPtr() {
	atomic.StoreUint32(&link.ptr, 0)
}

func (link *Linkables) Size() int {
	return int(atomic.LoadUint32(&link.size))
}

// SetupDataGathering
// sets the internal ptr of equal to the size
// now GetFile can be called until nil is returned and all values will be read correctly
// write is also set to false so no data can be modified
func (link *Linkables) SetupDataGathering() {
	link.mutex.Lock()
	link.ptr = link.size
	link.write = false
	link.mutex.Unlock()
}

func (link *Linkables) EnableWrite() {
	link.mutex.Lock()
	link.write = true
	link.mutex.Unlock()
}

func (link *Linkables) DisableWrite() {
	link.mutex.Lock()
	link.write = false
	link.mutex.Unlock()
}

func (link *Linkables) SetPtr(val uint32) {
	atomic.StoreUint32(&link.ptr, val)
}

func (link *Linkables) GetObjectFiles(outPath string, write bool) (objectFiles map[*pkg.ObjectFile]uint16, err error) {
	locations := make(map[uint16]uint16)
	objFiles := make(map[*pkg.ObjectFile]uint16)

	for _, file := range link.Files {
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
			objFile = assembler.AssembleAndWrite(string(file.Data), outPath, write)
		}

		// is location already used by another file
		if value, ok := locations[file.Location]; ok {
			objFiles[objFile] = file.Location + value
			locations[file.Location] = uint16(len(objFile.Code)) + value

		} else {
			objFiles[objFile] = file.Location
			locations[file.Location] = uint16(len(objFile.Code))
		}

	}

	return objFiles, nil
}
