package lexer

import (
	"fmt"

	"github.com/AirCraft009/mcc/internal/c"
	"github.com/AirCraft009/mcc/internal/helper"
)

func (lex *Lexer) Push(lexme string, Ttype c.TokenType) {
	lex.Size++
	curTok := lex.Tokens
	lex.Tokens = NewToken(lexme, Ttype)
	lex.Tokens.next = curTok
}

// Mod
//
// change the value of the topmost token
func (lex *Lexer) Mod(lexme string, Ttype c.TokenType) {
	lex.Tokens.lex = lexme
	lex.Tokens.TType = Ttype
}

func (lex *Lexer) Pop() *Token {
	lex.Size--
	token := lex.Tokens
	lex.Tokens = lex.Tokens.next
	return token
}

func (lex *Lexer) Peek() *Token {
	return lex.Tokens
}

func (lex *Lexer) Reset() {
	lex.Pos = 0
	lex.Size = 0
	lex.Tokens = nil
}

func (lex *Lexer) Parse(data string) error {
	lex.Reset()

	i := 0
	for i < len(data) {
		ch := standardizeSpaces(data[i])

		// STRING LIT
		if ch == '"' {
			start := i
			i++ // skip opening quote

			for i < len(data) {
				if data[i] == '\\' {
					i += 2 // skip escaped character
					continue
				}

				if data[i] == '"' {
					i++ // include closing quote
					break
				}

				i++
			}

			if i > len(data) || data[i-1] != '"' {
				return fmt.Errorf("unterminated string literal")
			}

			token := data[start:i]
			lex.AddToken(token, c.STRING, i)
			continue
		}

		// Skip whitespace
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			i++
			continue
		}

		//WORD(IDENT/KEYWORD)
		if helper.IsLetter(ch) {
			start := i
			i++

			for i < len(data) {
				c := data[i]
				if !helper.IsLetter(c) && !helper.IsDigit(c) {
					break
				}
				i++
			}

			token := data[start:i]

			if kw, ok := keywordMap[token]; ok {
				lex.AddToken(token, kw, i)
			} else {
				lex.AddToken(token, c.IDENT, i)
			}

			continue
		}

		// NUMBER
		if helper.IsDigit(ch) {
			start := i
			i++

			for i < len(data) && helper.IsDigit(data[i]) {
				i++
			}

			token := data[start:i]
			lex.AddToken(token, c.CONSTINT, i)
			continue
		}

		// DOUBLE CHAR OP
		if i+1 < len(data) {
			op := data[i : i+2]
			if t, ok := getOperator(op); ok {
				lex.AddToken(op, t, i+2)
				i += 2
				continue
			}
		}

		//SINGLE CHAR OP
		if t, ok := getOperator(string(ch)); ok {
			lex.AddToken(string(ch), t, i+1)
			i++
			continue
		}

		// DELIMITER
		if t, ok := getDelim(ch); ok {
			lex.AddToken(string(ch), t, i+1)
			i++
			continue
		}

		return fmt.Errorf("illegal character: %c", ch)
	}

	return nil
}

func (lex *Lexer) AddToken(token string, Ttype c.TokenType, newIndex int) {
	token = helper.StandardizeSpaceAmmount(token)
	lex.Push(token, Ttype)
	lex.Pos = newIndex
}

func (lex *Lexer) Output() {
	outStr := make([]string, lex.Size)
	token := lex.Tokens

	pos := lex.Size - 1
	for token != nil {
		outStr[pos] = fmt.Sprintf("%d=[%s: type %d]\n", pos, token.lex, token.TType)
		token = token.next
		pos--
	}

	fmt.Println(outStr)
}

func (lex *Lexer) ReturnTokens() []*Token {
	tokens := make([]*Token, 0, lex.Size)
	var pos int = int(lex.Size)
	for token := lex.Tokens; token != nil; token = lex.Tokens {
		pos--
		tokens[pos] = token
	}
	return tokens
}
