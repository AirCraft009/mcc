package main

import (
	"fmt"

	"github.com/AirCraft009/mcc/internal/c/lexer"
)

var code = "int main() {\n" +
	"	int a = 56;\n" +
	"	int b = 12;\n" +
	"	a += b;\n" +
	"	char buf[4] = \"afd\"" +
	"	int c5 = a + b;\n" +
	"	float *gd = 0;" +
	"	return c5;\n" +
	"}"

func main() {
	lexen := new(lexer.Lexer)
	/*
		err := lexen.Parse(code)
		if err != nil {
			fmt.Println("parse error:", err)
		}
		lexen.Output()
	*/

	err := lexen.Parse(code)
	if err != nil {
		fmt.Println("parse error:", err)
	}
	lexen.Output()
}
