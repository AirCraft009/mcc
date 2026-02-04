package mcc

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/AirCraft009/mcc/internal/helper"
)

//go:embed lib/stdlib/obj/.*obj
var embedFS embed.FS

type FSHelper struct {
	wdPath     string
	stdlibPath string
	EmbedFS    embed.FS
}

func InitFSHelper() *FSHelper {
	wdPath, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	stdlibPath := filepath.Join(helper.GetRootPath(), helper.StdLibLocationUse)

	return &FSHelper{wdPath, stdlibPath, embedFS}
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

	return FS.EmbedFS.ReadFile(filepath.Base(path))
}

func (FS *FSHelper) ResolveGlobal() {

}
