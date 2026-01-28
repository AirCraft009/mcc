package main

import (
	"mcc/Assembly-process/compiler"

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
		compiler.NoLinking(inputFile, *outPath)
		return
	}

	code, _ := compiler.NormalProcess(inputFile, false, false)
	err := os.WriteFile(*outPath, code, 0644)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("successfully wrote to: ", *outPath)
}
