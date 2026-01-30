package linker

import (
	"mcc/internal/assembler"
	helper2 "mcc/internal/helper"
	"mcc/pkg"
	"os"
	"path/filepath"
	"strconv"
)

// defineGlobalLookupTable
//
// parses the object files and puts all global labels
// (labels with an underscore at the beginning)
// and puts them into the GLT (Global Lookup Table)
// there they're related from symbol to global position
func defineGlobalLookupTable(objectFiles map[*pkg.ObjectFile]uint16) (globalLookupTable map[string]uint16) {
	globalLookupTable = make(map[string]uint16)

	for objFile, location := range objectFiles {
		for symbol, relAddr := range objFile.Symbols {
			if objFile.Globals[relAddr] {
				if _, ok := globalLookupTable[symbol]; ok {
					panic("Duplicate Lbl names: " + symbol + " location: " + strconv.Itoa(int(location)))
				}
				globalLookupTable[symbol] = location + relAddr
			}
		}
	}

	return globalLookupTable
}

func LinkModules(objectFiles map[*pkg.ObjectFile]uint16, debug, objectResolution bool) (code []byte, debugLocations map[uint16]string, err error) {
	//debug locations are only necesarry if the debugger is used can be discarded otherwise
	finalCode := make([]byte, helper2.MemorySize)
	if debug {
		debugLocations = make(map[uint16]string)
	}

	globalLookupTable := defineGlobalLookupTable(objectFiles)

	for objFile, location := range objectFiles {
		for _, relo := range objFile.Relocs {
			symbol, ok := objFile.Symbols[relo.Lbl]
			// is the label in the local scope ?
			if !ok {
				globalSymbol, k := globalLookupTable[relo.Lbl]
				// is the label global
				if !k {
					panic("Label not found: " + relo.Lbl)
				}
				symbol = globalSymbol
			} else {
				symbol += location
			}
			if debug {
				debugLocations[symbol] = relo.Lbl
			}
			hi, lo := helper2.EncodeAddr(symbol)
			objFile.Code[relo.Offset] = hi
			objFile.Code[relo.Offset+1] = lo
		}
		finalCode = helper2.ConcactSliceAtIndex(finalCode, objFile.Code, int(location))
		if objectResolution {
			//TODO: write the .obj next to the position of the .asm
			//pkg.SaveObjectFile(objFile, )
		}
	}
	return finalCode, debugLocations, nil
}

func CompileAndLinkFiles(files []string, originalLocations []uint16, outputPath string, debug bool) (code []byte, debugLocations map[uint16]string) {
	//for now this funcion will recomplile all files
	//It will take relative paths

	objFiles := make(map[*pkg.ObjectFile]uint16)
	locations := make(map[uint16]uint16)
	var obj *pkg.ObjectFile
	var err error

	for i := range len(files) {
		file, location := files[i], originalLocations[i]

		if filepath.Ext(file) == ".asm" {

			data, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}

			obj = assembler.Assemble(string(data), "", false)
		} else {
			// is already an object file or atleast not an asm file
			obj, err = pkg.ReadObjectFile(file)
			if err != nil {
				panic(err)
			}
		}

		// is location already used by another file
		if value, ok := locations[location]; ok {
			objFiles[obj] = location + value
			locations[location] = uint16(len(code)) + value

		} else {
			objFiles[obj] = location
			locations[location] = uint16(len(code))
		}

	}

	LinkedCode, debugLocations, err := LinkModules(objFiles, debug, false)

	if err != nil {
		panic(err.Error())
	}
	if outputPath == "" {
		outputPath = "a.exe"
	}

	err = os.WriteFile(outputPath, LinkedCode, 0644)
	if err != nil {
		panic(err.Error())
	}
	return LinkedCode, debugLocations
}
