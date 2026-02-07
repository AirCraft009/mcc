package main

import (
	"fmt"

	"github.com/AirCraft009/mcc/internal/fileHandling"
)

func findFileTest() error {
	// good Path
	includes, _, err := fileHandling.FindIncludes("examples/FindIncludes/main.asm")
	if err != nil {
		return err
	}
	fmt.Println("Good includes:", includes)

	_, _, err = fileHandling.FindIncludes("examples/FindIncludes/badExample.asm")
	if err == nil {
		return fmt.Errorf("was able to process badExample.asm\n")
	}

	return nil
}

func main() {
	fmt.Println("Testing: ")
	err := findFileTest()
	if err != nil {
		fmt.Println(err)
	}
}
