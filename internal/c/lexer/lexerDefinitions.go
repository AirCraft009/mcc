package lexer

import (
	"unicode"
)

type TokenType int

const (
	// Special
	ILLEGAL TokenType = iota
	EOF
	IDENT
	INT
	FLOAT
	DOUBLE
	CHAR

	// Keywords
	IF
	ELSE
	FOR
	WHILE
	DO
	SWITCH
	CASE
	DEFAULT
	BREAK
	CONTINUE
	RETURN
	CONST
	STRUCT
	ENUM
	STATIC
	EXTERN
	INLINE
	TRUE
	FALSE
	NULL

	// Operators
	ASSIGN    // =
	PLUS      // +
	MINUS     // -
	STAR      // *
	SLASH     // /
	PERCENT   // %
	INCREMENT // ++
	DECREMENT // --
	EQ        // ==
	NEQ       // !=
	LT        // <
	GT        // >
	LTE       // <=
	GTE       // >=
	AND       // &&
	OR        // ||
	NOT       // !
	BIT_AND   // &
	BIT_OR    // |
	BIT_XOR   // ^
	BIT_NOT   // ~
	SHL       // <<
	SHR       // >>

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	SEMICOLON // ;
	SPACE

	// forgot
	VOLATILE
	GOTO
	SHORT
	UNSIGNED
	SIGNED
	SIZEOF
	TYPEDEF
	UNION
	VOID
)

var keywordMap = map[string]TokenType{
	"break":    BREAK,
	"case":     CASE,
	"char":     CHAR,
	"const":    CONST,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"do":       DO,
	"double":   DOUBLE,
	"else":     ELSE,
	"enum":     ENUM,
	"extern":   EXTERN,
	"float":    FLOAT,
	"for":      FOR,
	"goto":     GOTO,
	"if":       IF,
	"inline":   INLINE,
	"int":      INT,
	"return":   RETURN,
	"short":    SHORT,
	"signed":   SIGNED,
	"sizeof":   SIZEOF,
	"static":   STATIC,
	"struct":   STRUCT,
	"switch":   SWITCH,
	"typedef":  TYPEDEF,
	"union":    UNION,
	"unsigned": UNSIGNED,
	"void":     VOID,
	"volatile": VOLATILE,
	"while":    WHILE,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
}

func getDelim(str string) (TokenType, bool) {
	switch str {
	case " ":
		return SPACE, true
	case "+":
		return PLUS, true
	case "-":
		return MINUS, true
	case "/":
		return SLASH, true
	case "*":
		return STAR, true
	case ",":
		return COMMA, true
	case ";":
		return SEMICOLON, true
	case "(":
		return LPAREN, true
	case ")":
		return RPAREN, true
	case "{":
		return LBRACE, true
	case "}":
		return RBRACE, true
	case "[":
		return LBRACKET, true
	case "]":
		return RBRACKET, true
	case "%":
		return PERCENT, true
	case "<":
		return LT, true
	case ">":
		return GT, true
	case "=":
		return ASSIGN, true
	default:
		return ILLEGAL, false
	}
}

func getOperator(str string) (TokenType, bool) {
	switch str {
	case "+":
		return PLUS, true
	case "-":
		return MINUS, true
	case "*":
		return STAR, true
	case "/":
		return SLASH, true
	case "<":
		return LT, true
	case ">":
		return GT, true
	case "=":
		return EQ, true
	case "++":
		return INCREMENT, true
	case "--":
		return DECREMENT, true
	case "!=":
		return NEQ, true
	case "<=":
		return LTE, true
	case ">=":
		return GTE, true
	case "&&":
		return AND, true
	case "||":
		return OR, true
	case "<<":
		return SHL, true
	case ">>":
		return SHR, true
	case "!":
		return NOT, true
	case "&":
		return BIT_AND, true
	case "|":
		return BIT_OR, true
	case "^":
		return BIT_XOR, true
	case "~":
		return BIT_NOT, true

	default:
		return ILLEGAL, false
	}
}

func isValidVarName(str string) bool {
	return unicode.IsLetter(rune(str[0]))
}

func isKeyword(str string) bool {
}
