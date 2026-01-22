package linker

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

const (
	ObjectF = iota + 1
	AsmF
	HeaderF
)

var acceptedFiletypes = map[string]byte{
	"asm": AsmF,
	"obj": ObjectF,
	"h":   HeaderF,
}

type Linkables struct {
	Files []*LinkFile
	ptr   uint32
	size  uint32
	write bool
	mutex sync.Mutex
}

type LinkFile struct {
	Path  string
	Data  []byte
	FileT uint8
}

func NewLinkables(size int) *Linkables {
	return &Linkables{make([]*LinkFile, size), 0, 0, true, sync.Mutex{}}
}

func (link *Linkables) AddFile(filePath string) (err error) {
	ext := filepath.Ext(filePath)
	if !(acceptedFiletypes[ext] == 0) {
		return errors.New("Unsupported file type: " + ext)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	linkF := &LinkFile{Path: filePath, Data: data, FileT: acceptedFiletypes[ext]}
	link.mutex.Lock()
	link.Files[link.ptr] = linkF
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
	link.ptr--
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
