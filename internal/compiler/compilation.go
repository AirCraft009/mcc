package compiler

import (
	"fmt"
	"strconv"

	"github.com/AirCraft009/mcc/internal/linker"
	preprocessor "github.com/AirCraft009/mcc/internal/pre-processor"
)

func NoLinking(inputFile, outPath string) {

	locations := make([]uint16, 1)
	includes := make([]string, 1)

	includes[0] = inputFile

	linker.IncludeHeaders(&includes, &locations)
	//fmt.Println(includes)
	link := linker.NewLinkables(len(includes))
	err := link.AddArraysMultiThreaded(includes, locations)

	if err != nil {
		panic(err.Error())
	}

	pre := preprocessor.NewPreProcesser()
	pre.Process(link)
	objs, err := link.GetObjectFiles(outPath, true)

	if err != nil {
		panic(err.Error())
	}
	// files still get written is just marked
	if len(objs) != 1 {
		panic("Expected 1 object got " + strconv.Itoa(len(objs)))
	}

	return
}

func NormalProcess(inputFile string, debug, resolution bool) ([]byte, map[uint16]string) {
	includes, locations, err := linker.FindIncludes(inputFile)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(includes)

	link := linker.NewLinkables(len(includes))
	fmt.Println(includes)
	err = link.AddArraysMultiThreaded(includes, locations)

	if err != nil {
		panic(err.Error())
	}

	pre := preprocessor.NewPreProcesser()
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
