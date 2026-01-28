package compiler

import (
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

	locations := make([]uint16, 0)
	includes := make([]string, 0)
	linker.IncludeBase(&includes, &locations)
	link := linker.NewLinkables(len(includes))

	err = link.AddArraysMultiThreaded(includes, locations)
	if err != nil {
		panic(err.Error())
	}
	pre := pre_processor.NewPreProcesser()
	pre.Process(link)
	_, err := link.GetObjectFiles(outPath, true)
	if err != nil {
		panic(err.Error())
	}
	return
}

func NormalProcess(inputFile string, debug, resolution bool) ([]byte, map[uint16]string) {
	includes, locations, err := linker.FindIncludes(inputFile)

	if err != nil {
		panic(err.Error())
	}
	link := linker.NewLinkables(len(includes))
	err = link.AddArraysMultiThreaded(includes, locations)
	if err != nil {
		panic(err.Error())
	}

	pre := pre_processor.NewPreProcesser()
	// define etc
	pre.Process(link)
	objs, err := link.GetObjectFiles("", false)
	if err != nil {
		panic(err.Error())
	}
	code, debugLabels, err := linker.LinkModules(objs, debug, resolution)
	if err != nil {
		panic(err.Error())
	}

	return code, debugLabels
}
