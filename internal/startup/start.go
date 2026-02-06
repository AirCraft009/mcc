package startup

import (
	"log"
	"strconv"

	"github.com/AirCraft009/mcc"
	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/internal/linker"
	preprocessor "github.com/AirCraft009/mcc/internal/pre-processor"
)

func NoLinking(inputFile, outPath string, logger *log.Logger) {

	locations := make([]uint16, 1)
	includes := make([]string, 1)

	includes[0] = inputFile
	fileSysHelper := mcc.InitFSHelper(logger)

	linker.IncludeHeaders(&includes, &locations)
	logger.Println("collected includes: ", includes)

	link := linker.NewLinkables(len(includes))
	err := link.AddArraysMultiThreaded(includes, locations, fileSysHelper)

	if err != nil {
		helper.FatalWrapper(logger, err.Error())
	}

	pre := preprocessor.NewPreProcesser()
	pre.Process(link)
	objs, _, err := link.GetObjectFiles(outPath, true, logger)

	if err != nil {
		helper.FatalWrapper(logger, err.Error())
	}
	// files still get written is just marked
	if len(objs) != 1 {
		helper.FatalWrapper(logger, "Expected 1 object got "+strconv.Itoa(len(objs)))
	}

	return
}

func NormalProcess(inputFile string, logger *log.Logger, debug, resolution bool) ([]byte, map[uint16]string) {
	logger.Println("starting linking assembly")
	logger.Println("finding includes")

	fileSysHelper := mcc.InitFSHelper(logger)
	includes, locations, err := linker.FindIncludes(inputFile, fileSysHelper)

	if err != nil {
		fileSysHelper.OutputVirtualFS()
		panic(err.Error())
	}
	logger.Println("collected includes: ", includes)

	link := linker.NewLinkables(len(includes))

	logger.Println("Adding files to linker")
	err = link.AddArraysMultiThreaded(includes, locations, fileSysHelper)

	if err != nil {
		panic(err.Error())
	}

	logger.Println("Successfully added files to linker")
	logger.Println("starting preprocessing")

	pre := preprocessor.NewPreProcesser()
	pre.Process(link)

	logger.Println("finished preprocessing")
	logger.Println("Assembling into Object Files")

	objs, data, err := link.GetObjectFiles("", false, logger)

	if err != nil {
		panic(err.Error())
	}

	logger.Println("Successfully made Object Files")
	logger.Println("Starting linking")

	code, debugLabels, err := linker.LinkModules(objs, data, debug, resolution, logger)
	if err != nil {
		panic(err.Error())
	}

	return code, debugLabels
}
