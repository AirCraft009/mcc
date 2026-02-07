package main

import (
	"flag"
	"fmt"

	"github.com/AirCraft009/mcc/internal/asm/startup"
	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

func main() {
	fk := helper.NewFlagKeeper()
	fk.Parse()
	// define flags

	flag.Parse()
	// init the main logger
	myLogger := helper.InitLogger(fk.Supress, fk.Verbose)

	myLogger.Printf("Compiling %s to %s\n", fk.InputPath, fk.OutPath)
	if fk.NoLink {
		startup.NoLinking(fk.InputPath, fk.OutPath, myLogger)
		return
	}

	code, debugSymbols := startup.NormalProcess(fk.InputPath, myLogger, fk.Debug, fk.Resolution)
	err := pkg.WriteMxBinary(fk.OutPath, code, debugSymbols, myLogger, fk.Debug)
	if err != nil {
		helper.FatalWrapper(myLogger, err.Error())
	}

	fmt.Println("successfully wrote to: ", fk.OutPath)
	return
}
