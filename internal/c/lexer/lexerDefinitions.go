package lexer

import (
	"errors"
	"strconv"
	"unicode"

	def "github.com/AirCraft009/mcc/internal/c"
)

type Lexer struct {
	Size   uint64
	Pos    int
	Tokens *Token
}

type Token struct {
	TType def.TokenType
	Lex   string
	next  *Token
}

func NewToken(lex string, ttype def.TokenType) *Token {
	return &Token{
		Lex:   lex,
		TType: ttype,
		next:  nil,
	}
}

var keywordMap = map[string]def.TokenType{
	"break":    def.BREAK,
	"case":     def.CASE,
	"char":     def.CHAR,
	"const":    def.CONST,
	"continue": def.CONTINUE,
	"default":  def.DEFAULT,
	"do":       def.DO,
	"double":   def.DOUBLE,
	"else":     def.ELSE,
	"enum":     def.ENUM,
	"extern":   def.EXTERN,
	"float":    def.FLOAT,
	"for":      def.FOR,
	"goto":     def.GOTO,
	"if":       def.IF,
	"inline":   def.INLINE,
	"int":      def.INT,
	"return":   def.RETURN,
	"short":    def.SHORT,
	"signed":   def.SIGNED,
	"sizeof":   def.SIZEOF,
	"static":   def.STATIC,
	"struct":   def.STRUCT,
	"switch":   def.SWITCH,
	"typedef":  def.TYPEDEF,
	"union":    def.UNION,
	"unsigned": def.UNSIGNED,
	"void":     def.VOID,
	"volatile": def.VOLATILE,
	"while":    def.WHILE,
	"true":     def.TRUE,
	"false":    def.FALSE,
	"null":     def.NULL,
}

func GetDelim(ch byte) (def.TokenType, bool) {
	switch ch {
	case ' ':
		return def.SPACE, true
	case '+':
		return def.PLUS, true
	case '-':
		return def.MINUS, true
	case '/':
		return def.SLASH, true
	case '*':
		return def.STAR, true
	case ',':
		return def.COMMA, true
	case ';':
		return def.SEMICOLON, true
	case '(':
		return def.LPAREN, true
	case ')':
		return def.RPAREN, true
	case '{':
		return def.LBRACE, true
	case '}':
		return def.RBRACE, true
	case '[':
		return def.LBRACKET, true
	case ']':
		return def.RBRACKET, true
	case '%':
		return def.PERCENT, true
	case '<':
		return def.LT, true
	case '>':
		return def.GT, true
	case '=':
		return def.ASSIGN, true
	default:
		return def.ILLEGAL, false
	}
}

func GetOperator(str string) (def.TokenType, bool) {
	switch str {
	case "+":
		return def.PLUS, true
	case "-":
		return def.MINUS, true
	case "*":
		return def.STAR, true
	case "/":
		return def.SLASH, true
	case "<":
		return def.LT, true
	case ">":
		return def.GT, true
	case "=":
		return def.EQ, true
	case "+=":
		return def.PLUSEQ, true
	case "-=":
		return def.MINUSEQ, true
	case "*=":
		return def.MULEQ, true
	case "/=":
		return def.DIVEQ, true
	case "++":
		return def.INCREMENT, true
	case "--":
		return def.DECREMENT, true
	case "!=":
		return def.NEQ, true
	case "<=":
		return def.LTE, true
	case ">=":
		return def.GTE, true
	case "&&":
		return def.AND, true
	case "||":
		return def.OR, true
	case "<<":
		return def.SHL, true
	case ">>":
		return def.SHR, true
	case "!":
		return def.NOT, true
	case "&":
		return def.BIT_AND, true
	case "|":
		return def.BIT_OR, true
	case "^":
		return def.BIT_XOR, true
	case "~":
		return def.BIT_NOT, true

	default:
		return def.ILLEGAL, false
	}
}

func isValidIdentifier(str string) bool {
	if len(str) == 0 {
		return false
	}
	return unicode.IsLetter(rune(str[0]))
}

func standardizeSpaces(r byte) byte {
	if unicode.IsSpace(rune(r)) {
		return ' '
	}

	return r
}

func checkIdentifier(s string) (def.TokenType, error) {
	if len(s) == 0 {
		return def.ILLEGAL, errors.New("identifier is empty")
	}
	fc := s[0] - '0'
	if fc < 10 {
		_, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return def.CONSTINT, nil
	}

	return def.IDENT, nil
}
