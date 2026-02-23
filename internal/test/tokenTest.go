package main

import (
	"fmt"

	"github.com/AirCraft009/mcc/internal/c/Parser"
	"github.com/AirCraft009/mcc/internal/c/lexer"
)

var code = "int main(int input) {\n" +
	"	int a = 56;\n" +
	"	int b = 12;\n" +
	"	add(input, a)\n" +
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

	tokens := lexen.ReturnTokens()
	parse := Parser.NewParser(tokens)
	statements, err := parse.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", statements)
}
