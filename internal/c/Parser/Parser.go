package Parser

import (
	"fmt"
	"strconv"

	"github.com/AirCraft009/mcc/internal/c"
	"github.com/AirCraft009/mcc/internal/c/lexer"
)

//
// AST INTERFACES
//

type Statements interface {
	String() string
}

//
// AST NODES
//

type FuncDecl struct {
	FuncName string
	FuncArgs []Statements
	RetType  string
	Body     *Body
}

func (f *FuncDecl) String() string {
	return fmt.Sprintf("FuncDecl(%s %s %v %s)",
		f.RetType,
		f.FuncName,
		f.FuncArgs,
		f.Body)
}

type Body struct {
	Blocks []Statements
}

func (b *Body) String() string {
	return fmt.Sprintf("{ %v }", b.Blocks)
}

type ArgVariable struct {
	varType string
	name    string
}

func (v *ArgVariable) String() string {
	return fmt.Sprintf("%s %s", v.varType, v.name)
}

type methodCall struct {
	name string
	args []Statements
}

func (m *methodCall) String() string {
	return fmt.Sprintf("%s(%v)", m.name, m.args)
}

type refVar struct {
	name string
}

func (r *refVar) String() string {
	return r.name
}

type returnStmt struct {
	retExpr Statements
}

func (r *returnStmt) String() string {
	return fmt.Sprintf("return %s;", r.retExpr)
}

type OpExpr struct {
	Left  Statements
	Right Statements
	Op    c.TokenType
}

func (o *OpExpr) String() string {
	return fmt.Sprintf("(%s %v %s)", o.Left, o.Op, o.Right)
}

type IntegerLit struct {
	Value int64
}

func (i *IntegerLit) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type StringLit struct {
	Value string
}

func (s *StringLit) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

//
// PARSER
//

type Parser struct {
	tokens []*lexer.Token
	pos    int
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

//
// TOKEN HELPERS
//

func (p *Parser) current() *lexer.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() *lexer.Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) match(tt c.TokenType) bool {
	if p.current() != nil && p.current().TType == tt {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) expect(tt c.TokenType) (*lexer.Token, error) {

	tok := p.current()

	if tok == nil {
		return nil, fmt.Errorf("unexpected EOF, expected %v", tt)
	}

	if tok.TType != tt {
		return nil, fmt.Errorf("expected %v got %v", tt, tok.TType)
	}

	p.pos++
	return tok, nil
}

func (p *Parser) expectType() (*lexer.Token, error) {

	tok := p.current()

	if tok == nil {
		return nil, fmt.Errorf("unexpected EOF, expected type")
	}

	if tok.TType >= c.ASSIGN && tok.TType >= c.SPACE {
		return nil, fmt.Errorf("expected type got delim: %v", tok.TType)
	}

	// just takes anything

	p.pos++
	return tok, nil
}

//
// TOP LEVEL
//

func (p *Parser) Parse() ([]Statements, error) {

	var nodes []Statements

	for p.current() != nil {

		fn, err := p.parseFuncDecl()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, fn)
	}

	return nodes, nil
}

//
// FUNCTION DECL
//

func (p *Parser) parseFuncDecl() (*FuncDecl, error) {

	retTypeTok, err := p.expectType()
	if err != nil {
		return nil, err
	}

	nameTok, err := p.expect(c.IDENT)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(c.LPAREN)
	if err != nil {
		return nil, err
	}

	fn := &FuncDecl{
		RetType:  retTypeTok.Lex,
		FuncName: nameTok.Lex,
		FuncArgs: []Statements{},
	}

	for !p.match(c.RPAREN) {

		typeTok, err := p.expectType()
		if err != nil {
			return nil, err
		}

		nameTok, err := p.expect(c.IDENT)
		if err != nil {
			return nil, err
		}

		fn.FuncArgs = append(fn.FuncArgs,
			&ArgVariable{
				varType: typeTok.Lex,
				name:    nameTok.Lex,
			})

		p.match(c.COMMA)
	}

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}

	fn.Body = body

	return fn, nil
}

//
// BODY
//

func (p *Parser) parseBody() (*Body, error) {

	_, err := p.expect(c.LBRACE)
	if err != nil {
		return nil, err
	}

	body := &Body{
		Blocks: []Statements{},
	}

	for !p.match(c.RBRACE) {

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		body.Blocks = append(body.Blocks, stmt)
	}

	return body, nil
}

//
// STATEMENT
//

func (p *Parser) parseStatement() (Statements, error) {

	tok := p.current()

	if tok == nil {
		return nil, fmt.Errorf("unexpected EOF")
	}

	switch tok.TType {

	case c.RETURN:
		return p.parseReturn()

	case c.IDENT:

		// ensure this is actually a function call
		if p.pos+1 >= len(p.tokens) || p.tokens[p.pos+1].TType != c.LPAREN {
			return nil, fmt.Errorf(
				"unexpected identifier '%s', expected function call",
				tok.Lex,
			)
		}

		call, err := p.parseMethodCall()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(c.SEMICOLON)
		if err != nil {
			return nil, err
		}

		return call, nil

	case c.INT, c.SHORT, c.CHAR:

	}
	return nil, fmt.Errorf("unexpected token %v", tok.TType)
}

//
// RETURN
//

func (p *Parser) parseReturn() (*returnStmt, error) {

	p.advance()

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(c.SEMICOLON)
	if err != nil {
		return nil, err
	}

	return &returnStmt{
		retExpr: expr,
	}, nil
}

//
// METHOD CALL
//

func (p *Parser) parseMethodCall() (*methodCall, error) {

	nameTok, err := p.expect(c.IDENT)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(c.LPAREN)
	if err != nil {
		return nil, err
	}

	call := &methodCall{
		name: nameTok.Lex,
		args: []Statements{},
	}

	for !p.match(c.RPAREN) {

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		call.args = append(call.args, expr)

		p.match(c.COMMA)
	}

	return call, nil
}

//
// EXPRESSIONS
//

func (p *Parser) parseExpression() (Statements, error) {
	return p.parseBinary()
}

func (p *Parser) parseBinary() (Statements, error) {

	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {

		tok := p.current()
		if tok == nil {
			break
		}

		switch tok.TType {

		case c.PLUS, c.MINUS, c.STAR, c.SLASH:

			op := tok.TType
			p.advance()

			right, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}

			left = &OpExpr{
				Left:  left,
				Right: right,
				Op:    op,
			}

		default:
			return left, nil
		}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (Statements, error) {

	tok := p.current()

	switch tok.TType {

	case c.INT:

		p.advance()

		val, err := strconv.ParseInt(tok.Lex, 10, 64)
		if err != nil {
			return nil, err
		}

		return &IntegerLit{
			Value: val,
		}, nil

	case c.STRING:

		p.advance()

		return &StringLit{
			Value: tok.Lex,
		}, nil

	case c.IDENT:

		if p.tokens[p.pos+1].TType == c.LPAREN {
			return p.parseMethodCall()
		}

		p.advance()

		return &refVar{
			name: tok.Lex,
		}, nil

	case c.LPAREN:

		p.advance()

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(c.RPAREN)
		return expr, err
	}

	return nil, fmt.Errorf("unexpected token %v", tok.TType)
}
