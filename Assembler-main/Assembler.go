package main

import (
	"mcc/Assembly-process/assembler"
	"mcc/Assembly-process/linker"
	"strings"

	flag "github.com/spf13/pflag"
)
import (
	"fmt"
	"os"
)

func main() {
	// define flags
	outPath := flag.String("o", "a.bin", "output file")
	noLink := flag.Bool("no_link", false, "do not use linker")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Printf("Usage: ./mcc inputFile -o outputfile\n")
		return
	}

	inputFile := os.Args[1]
	if *noLink {
		file, err := os.ReadFile(inputFile)
		if err != nil {
			panic(err.Error())
		}
		assembler.Assemble(string(file), strings.Clone(*outPath), true)
		return
	}
	includes, locations, err := linker.FindIncludes(inputFile)
	if err != nil {
		panic(err)
	}

	linker.CompileAndLinkFiles(includes, locations, strings.Clone(*outPath), false)
}
