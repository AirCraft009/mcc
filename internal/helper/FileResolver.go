package helper

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
)

//go:embed lib/stdlib/obj/*.obj
var EmbedFS embed.FS

type FSHelper struct {
	wdPath     string
	stdlibPath string
	EmbedFS    embed.FS
}

func initFSHelper() *FSHelper {
	wdPath, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	stdlibPath := filepath.Join(GetRootPath(), StdLibLocationUse)

	return &FSHelper{wdPath, stdlibPath, EmbedFS}
}

func (FS *FSHelper) ResolveReadFile(path string) ([]byte, error) {
	wdPathFull := filepath.Clean(filepath.Join(FS.wdPath, path))
	if fileData, err := os.ReadFile(wdPathFull); err == nil {
		return fileData, nil
	}

	stdlibFullPath := filepath.Clean(filepath.Join(FS.stdlibPath, path))
	if fileData, err := os.ReadFile(stdlibFullPath); err == nil {
		return fileData, nil
	}

	if fileData, err := FS.EmbedFS.ReadFile(filepath.Base(path)); err == nil {
		return fileData, nil
	}

	return nil, errors.New("No File Found while resolving: " + path)
}

func (FS *FSHelper) ResolveGlobal() {

}
