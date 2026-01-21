package main

import (
	"fmt"
	"mcc/Assembly-process/linker"
)

func findFileTest() error {
	// good Path
	includes, err := linker.FindIncludes("examples/FindIncludes/main.asm")
	if err != nil {
		return err
	}
	fmt.Println("Good includes:", includes)

	_, err = linker.FindIncludes("examples/FindIncludes/badExample.asm")
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
