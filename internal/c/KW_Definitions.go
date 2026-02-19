package c

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
	// Constant
	STRING
	CONSTINT

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
	PLUSEQ
	MINUSEQ
	DIVEQ
	MULEQ

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
