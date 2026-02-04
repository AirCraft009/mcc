package main

import (
	"flag"
	"fmt"

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

	flag.Parse()
	if *verbose {
		fmt.Printf("Compiling %s to %s\n", *inputFile, *outPath)
	}

	if *noLink {
		startup.NoLinking(*inputFile, *outPath, *verbose)
		return
	}

	code, debugSymbols := startup.NormalProcess(*inputFile, *debug, *resolution, *verbose)
	err := pkg.WriteMxBinary(*outPath, code, debugSymbols, *debug)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("successfully wrote to: ", *outPath)
	return
}
