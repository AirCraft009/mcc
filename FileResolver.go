package mcc

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/internal/mcc-constants"
)

//go:embed lib/stdlib/obj/*
var embedFS embed.FS

// FSHelper
//
// Contains information to Resolve filepaths
// For mcc
type FSHelper struct {
	wdPath      string
	Base        string
	stdlibPath  string
	EmbedFsBase string
	EmbedFS     embed.FS
}

// InitFSHelper
//
// Create a new FSHelper
func InitFSHelper() *FSHelper {
	wdPath, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	stdlibPath := filepath.Join(helper.GetRootPath(), mcc_constants.StdLibLocationUse)

	return &FSHelper{wdPath, helper.GetRootPath(), stdlibPath, mcc_constants.StdLibLocationUse, embedFS}
}

func (FS *FSHelper) OutputVirtualFS() {
	fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println("File:", path)
		return nil
	})
}

// ResolveReadFile
//
// replaces normal os.ReadFile()
//
// This handles reading files in this order:
//
// - 1. Relative filepaths
// - 2. stdlib Paths at real location
// - 3. stdlib Paths embedded
// - 4. Absolute path (os.ReadFile())
func (FS *FSHelper) ResolveReadFile(path string) ([]byte, error) {

	// relative path shadows all
	wdPathFull := filepath.Clean(filepath.Join(FS.wdPath, path))
	if fileData, err := os.ReadFile(wdPathFull); err == nil {
		return fileData, nil
	}

	// is in the stdlib
	if filepath.Clean(filepath.Dir(path)) == mcc_constants.StdLibLocationSignifier {
		// The files are actually present at the location
		stdlibFullPath := filepath.Clean(filepath.Join(FS.stdlibPath, filepath.Base(path)))
		if fileData, err := os.ReadFile(stdlibFullPath); err == nil {
			return fileData, nil
		}

		// The file is embedded
		return FS.EmbedFS.ReadFile(filepath.Join(FS.EmbedFsBase, filepath.Base(path)))
	}

	// real absolute path
	return os.ReadFile(path)
}

func (FS *FSHelper) ResolveGlobal() {

}
