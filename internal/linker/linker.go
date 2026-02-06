package linker

import (
	"log"
	"strconv"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

// defineGlobalLookupTable
//
// parses the object files and puts all global labels
// (labels with an underscore at the beginning)
// and puts them into the GLT (Global Lookup Table)
// there they're related from symbol to global position
func defineGlobalLookupTable(objectFiles map[*pkg.ObjectFile]uint16, logger *log.Logger) (globalLookupTable map[string]uint16) {
	globalLookupTable = make(map[string]uint16)

	for objFile, location := range objectFiles {
		for symbol, relAddr := range objFile.Symbols {
			if objFile.Globals[relAddr] {
				if _, ok := globalLookupTable[symbol]; ok {
					helper.FatalWrapper(logger, "Duplicate Lbl names: "+symbol+" location: "+strconv.Itoa(int(location)))
				}
				globalLookupTable[symbol] = location + relAddr
			}
		}
	}

	return globalLookupTable
}

func LinkModules(objectFiles map[*pkg.ObjectFile]uint16, Datasection []byte, debug, objectResolution bool, logger *log.Logger) (code []byte, debugLocations map[uint16]string, err error) {
	//debug locations are only necesarry if the debugger is used can be discarded otherwise
	finalCode := make([]byte, pkg.MemorySize)
	if debug {
		debugLocations = make(map[uint16]string)
	}

	globalLookupTable := defineGlobalLookupTable(objectFiles, logger)

	for objFile, location := range objectFiles {
		logger.Printf("Linking to %d\n", location)
		for _, relo := range objFile.Relocs {
			symbol, ok := objFile.Symbols[relo.Lbl]

			// is the label in the local scope ?
			if !ok {
				globalSymbol, k := globalLookupTable[relo.Lbl]
				// is the label global
				if !k {
					helper.FatalWrapper(logger, "Lable not found: "+relo.Lbl)
				}
				symbol = globalSymbol
			} else if !relo.Data {
				symbol += location
			}

			if debug {
				debugLocations[symbol] = relo.Lbl
			}
			logger.Printf("Linking label: %s to addr: %d\n", relo.Lbl, symbol)

			hi, lo := helper.EncodeAddr(uint16(int32(symbol) + relo.Offset))
			objFile.Code[relo.InFileOffset] = hi
			objFile.Code[relo.InFileOffset+1] = lo
		}

		finalCode = helper.ConcactSliceAtIndex(finalCode, objFile.Code, int(location))
		if objectResolution {
			//TODO: write the .obj next to the position of the .asm
			//pkg.SaveObjectFile(objFile, )
		}

		if (objFile.Relocs == nil || len(objFile.Relocs) == 0) && debug {

			for name, relAddr := range objFile.Symbols {
				debugLocations[relAddr+location] = name
			}
		}

	}

	//fmt.Println("final code", Datasection)
	copy(finalCode[pkg.DataStart:pkg.DataEnd+1], Datasection)
	return finalCode, debugLocations, nil
}
