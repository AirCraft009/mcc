package startup

import (
	"fmt"
	"strconv"

	"github.com/AirCraft009/mcc"
	"github.com/AirCraft009/mcc/internal/linker"
	preprocessor "github.com/AirCraft009/mcc/internal/pre-processor"
)

func NoLinking(inputFile, outPath string, verbose bool) {

	locations := make([]uint16, 1)
	includes := make([]string, 1)

	includes[0] = inputFile
	fileSysHelper := mcc.InitFSHelper()

	linker.IncludeHeaders(&includes, &locations)
	if verbose {
		fmt.Println("collected includes: ", includes)
	}
	link := linker.NewLinkables(len(includes))
	err := link.AddArraysMultiThreaded(includes, locations, fileSysHelper)

	if err != nil {
		panic(err.Error())
	}

	pre := preprocessor.NewPreProcesser()
	pre.Process(link)
	objs, _, err := link.GetObjectFiles(outPath, true, verbose)

	if err != nil {
		panic(err.Error())
	}
	// files still get written is just marked
	if len(objs) != 1 {
		panic("Expected 1 object got " + strconv.Itoa(len(objs)))
	}

	return
}

func NormalProcess(inputFile string, debug, resolution, verbose bool) ([]byte, map[uint16]string) {
	if verbose {
		fmt.Println("starting linking assembly")
		fmt.Println("finding includes")
	}
	fileSysHelper := mcc.InitFSHelper()
	includes, locations, err := linker.FindIncludes(inputFile, fileSysHelper)

	if err != nil {
		fileSysHelper.OutputVirtualFS()
		panic(err.Error())
	}
	if verbose {
		fmt.Println("collected includes: ", includes)
	}

	link := linker.NewLinkables(len(includes))
	if verbose {
		fmt.Println("Adding files to linker")
	}
	err = link.AddArraysMultiThreaded(includes, locations, fileSysHelper)

	if err != nil {
		panic(err.Error())
	}

	if verbose {
		fmt.Println("Successfully added files to linker")
		fmt.Println("starting preprocessing")
	}

	pre := preprocessor.NewPreProcesser()
	pre.Process(link)

	if verbose {
		fmt.Println("finished preprocessing")
		fmt.Println("Assembling into Object Files")
	}
	objs, data, err := link.GetObjectFiles("", false, verbose)

	if err != nil {
		panic(err.Error())
	}
	if verbose {
		fmt.Println("Successfully made Object Files")
		fmt.Println("Starting linking")
	}
	code, debugLabels, err := linker.LinkModules(objs, data, debug, resolution, verbose)
	if err != nil {
		panic(err.Error())
	}

	return code, debugLabels
}
