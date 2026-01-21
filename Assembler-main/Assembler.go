package main

import (
	"mcc/Assembly-process/linker"

	flag "github.com/spf13/pflag"
)
import (
	"fmt"
	"os"
)

func main() {
	// define flags
	_ = flag.String("o", "out.bin", "output file")
	if len(os.Args) < 2 {
		fmt.Printf("Usage: ./mcc inputFile -o outputfile\n")
		return
	}

	inputFile := os.Args[1]

	_, err := linker.FindIncludes(inputFile)
	if err != nil {
		panic(err)
	}
}
