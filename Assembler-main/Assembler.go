package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/internal/startup"
	"github.com/AirCraft009/mcc/pkg"
)

func main() {
	// define flags
	inputFile := flag.String("i", "", "input file")
	outPath := flag.String("o", "a.bin", "output file")
	noLink := flag.Bool("n", false, "do not use linker\n overrides debug and res because no full file is created")
	debug := flag.Bool("debug", false, "creates debug symbols")
	resolution := flag.Bool("res", false, "creates the object files at in the dir next to eachother")
	verbose := flag.Bool("v", false, "verbose output")
	supress := flag.Bool("s", false, "doesn't write to a log file at all")
	flag.Parse()
	// init the main logger
	myLogger := initLogger(*supress, *verbose)

	// init the logger for warnings (always to stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("MCC-WARN: ")

	myLogger.Printf("Compiling %s to %s\n", *inputFile, *outPath)
	if *noLink {
		startup.NoLinking(*inputFile, *outPath, myLogger)
		return
	}

	code, debugSymbols := startup.NormalProcess(*inputFile, myLogger, *debug, *resolution)
	err := pkg.WriteMxBinary(*outPath, code, debugSymbols, myLogger, *debug)
	if err != nil {
		helper.FatalWrapper(myLogger, err.Error())
	}

	fmt.Println("successfully wrote to: ", *outPath)
	return
}

func initLogger(supress, verbose bool) *log.Logger {
	var myLogger *log.Logger

	if !supress {
		if !verbose {
			create, err := os.Create("Mcc-Logger.log")
			if err != nil {
				return log.New(io.Discard, "", 0)
			}
			fmt.Printf("writing logs to: %s\n", create.Name())
			myLogger = log.New(create, "Mcc-assmbler:", log.LstdFlags|log.Lshortfile)

		} else {
			myLogger = log.New(os.Stderr, "Mcc-assembler:", log.LstdFlags)
		}

	} else {
		myLogger = log.New(io.Discard, "", 0)
	}
	return myLogger
}
