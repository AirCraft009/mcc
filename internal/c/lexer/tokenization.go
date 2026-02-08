package lexer

import (
	"fmt"
	"strings"
)

func (lex *Lexer) Push(lexme string, Ttype TokenType) {
	lex.size++
	curTok := lex.Tokens
	lex.Tokens = NewToken(lexme, Ttype)
	lex.Tokens.next = curTok
}

func (lex *Lexer) Pop() *Token {
	lex.size--
	token := lex.Tokens
	lex.Tokens = lex.Tokens.next
	return token
}

func (lex *Lexer) Peek() *Token {
	return lex.Tokens
}

func (lex *Lexer) Parse(data string) {
	// I wrote this without any input it might be bad, but it's working

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		for _, field := range fields {
			buildStr := []byte(field)
			var lastIndex int
			for i := range buildStr {
				char := string(buildStr[i])
				delim, ok := GetDelim(char)

				if ok {
					tokenStr := string(buildStr[lastIndex:i])
					if tokenStr == "" {
						continue
					}
					tt := getTokenType(tokenStr)

					lex.Push(tokenStr, tt)
					lex.Push(char, delim)
					// don't push the delim next time
					lastIndex = i + 1
				}

			}
			if lastIndex < len(buildStr) {
				tokenStr := string(buildStr[lastIndex:])
				if tokenStr == "" {
					continue
				}
				tt := getTokenType(tokenStr)
				lex.Push(tokenStr, tt)
			}
		}
	}
}

func (lex *Lexer) Output() {
	outStr := make([]string, lex.size)
	token := lex.Tokens

	pos := lex.size - 1
	for token != nil {
		outStr[pos] = fmt.Sprintf("[%s: type %d]\n", token.lex, token.TType)
		token = token.next
		pos--
	}

	fmt.Println(outStr)
}
