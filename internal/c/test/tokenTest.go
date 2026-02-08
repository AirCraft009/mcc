package main

import "github.com/AirCraft009/mcc/internal/c/lexer"

var code = "int main() {\n" +
	"	int a = 56;\n" +
	"	int b = 12;\n" +
	"	int c = a + b;\n" +
	"	return c;\n" +
	"}"

func main() {
	lexen := new(lexer.Lexer)

	lexen.Parse(code)
	lexen.Output()
}
