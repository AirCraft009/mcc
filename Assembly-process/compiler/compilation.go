package compiler

import (
	"fmt"
	"mcc/Assembly-process/assembler"
	"mcc/Assembly-process/linker"
	pre_processor "mcc/Assembly-process/pre-processor"
	"os"
	"strings"
)

func NoLinking(inputFile, outPath string) {
	file, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err.Error())
	}
	assembler.Assemble(string(file), strings.Clone(outPath), true)
	return
}

func NormalProcess(inputFile string, debug, resolution bool) ([]byte, map[uint16]string) {
	includes, locations, err := linker.FindIncludes(inputFile)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("includes:", includes)
	link := linker.NewLinkables(len(includes))
	fmt.Println("linker:", link)
	err = link.AddArrays(includes, locations)
	fmt.Println("added arrays: ", link)
	if err != nil {
		panic(err.Error())
	}

	pre := pre_processor.NewPreProcesser()
	fmt.Println("pre processor:", pre)
	// define etc
	pre.Process(link)
	fmt.Println("pre processed")

	objs, err := link.GetObjectFiles()
	fmt.Println("got objects")
	if err != nil {
		panic(err.Error())
	}
	code, debugLabels, err := linker.LinkModules(objs, debug, resolution)
	fmt.Println("final code:", code)
	if err != nil {
		panic(err.Error())
	}

	return code, debugLabels
}
