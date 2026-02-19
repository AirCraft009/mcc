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

	decl := new(Parser.FuncDecl)
	decl.ParseSelf(tokens)
	fmt.Println(decl.FuncName)
	fmt.Println(decl.FuncArgs)
	fmt.Println(decl.RetType)

	fmt.Println("parsing body")
	fmt.Println(tokens[6].Lex)

	bdy := new(Parser.Body)
	bdy.ParseSelf(lexen.ReturnTokens()[6:])
	fmt.Println(bdy.Blocks)
}
