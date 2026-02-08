package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AirCraft009/mcc/internal/asm/startup"
	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

func main() {
	fk := new(helper.FlagKeeper)
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flagSet.StringVar(&fk.OutPath, "o", "a.bin", "output file")
	flagSet.BoolVar(&fk.NoLink, "n", false, "do not use linker\n\toverrides debug and res because no full file is created")
	flagSet.BoolVar(&fk.Debug, "debug", false, "creates debug symbols")
	flagSet.BoolVar(&fk.Resolution, "res", false, "creates the object files at in the dir next to eachother")
	flagSet.BoolVar(&fk.Verbose, "v", false, "verbose output")
	flagSet.BoolVar(&fk.Log, "log", false, "writes to a log file")
	flagSet.BoolVar(&fk.Supress, "s", true, "doesn't write to a log file at all - log output to stderr")

	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		panic(err.Error())
	}
	fk.InputPath = os.Args[1]

	// init the main logger
	myLogger := helper.InitLogger(fk.Supress, fk.Verbose, fk.Log)

	myLogger.Printf("Compiling %s to %s\n", fk.InputPath, fk.OutPath)
	if fk.NoLink {
		startup.NoLinking(fk.InputPath, fk.OutPath, myLogger)
		return
	}

	code, debugSymbols := startup.NormalProcess(fk.InputPath, myLogger, fk.Debug, fk.Resolution)
	err = pkg.WriteMxBinary(fk.OutPath, code, debugSymbols, myLogger, fk.Debug)
	if err != nil {
		helper.FatalWrapper(myLogger, err.Error())
	}

	fmt.Println("successfully wrote to: ", fk.OutPath)
	return
}
