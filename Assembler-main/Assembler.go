package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AirCraft009/mcc/internal/compiler"
	"github.com/AirCraft009/mcc/pkg"
)

func main() {
	// define flags
	outPath := flag.String("o", "a.bin", "output file")
	noLink := flag.Bool("no_link", false, "do not use linker\n overrides debug and res because no full file is created")
	debug := flag.Bool("debug", false, "creates debug symbols")
	resolution := flag.Bool("res", false, "creates the object files at in the dir next to eachother")

	flag.Parse()

	if len(os.Args) < 2 {
		_ = fmt.Errorf("mcc: {No Input File specified.\n}")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	if *noLink {
		compiler.NoLinking(inputFile, *outPath)
		return
	}

	code, debugSymbols := compiler.NormalProcess(inputFile, *debug, *resolution)
	err := pkg.WriteMxBinary(*outPath, code, debugSymbols, *debug)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("successfully wrote to: ", *outPath)
	return
}
