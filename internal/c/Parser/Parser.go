package Parser

import (
	"errors"
	"fmt"

	"github.com/AirCraft009/mcc/internal/c"
	"github.com/AirCraft009/mcc/internal/c/lexer"
)

// Statements
// Abstract Syntax Tree
type Statements interface {
	parseSelf(data []*lexer.Token) error
	String() string
}

type FuncDecl struct {
	FuncName string
	FuncArgs []Statements
	RetType  string
}

// ParseSelf
// after a type ident(
// but data is passed starting at -> type
func (f *FuncDecl) ParseSelf(data []*lexer.Token) error {
	// following is either
	// a type if the func has params
	// a RParen if it doesn't
	f.FuncArgs = make([]Statements, 0, 3)

	f.RetType = data[0].Lex
	f.FuncName = data[1].Lex

	// cutoff the RPAREN
	data = data[3:]
	for i := 0; i < len(data); i++ {
		token := data[i]
		if token.TType == c.RPAREN {
			break
		}

		// delim between two vars
		if token.TType == c.COMMA {
			continue
		}

		// handle a var
		// first char hast to be a datatype
		f.FuncArgs = append(
			f.FuncArgs, &ArgVariable{
				varType: data[i].Lex,
				name:    data[i+1].Lex,
			})
		// add one extra because I used i + 1 for the name
		i++
	}
	return nil
}

func (f *FuncDecl) String() string {
	return fmt.Sprintf("%s %s(%s)", f.FuncName, f.RetType, f.FuncArgs)
}

/*

SNIPPET:

	data = data[i+1:]
	if data[0].TType != c.LBRACE {
		// function definitions do the same
		// may change in the future

		return errors.New("function declaration without left brace")
	}
*/

type Body struct {
	Blocks []Statements
}

func (b *Body) ParseSelf(data []*lexer.Token) error {
	b.Blocks = make([]Statements, 0, 10)

	for i := 0; i < len(data); i++ {
		token := data[i]

		if token.TType == c.RBRACE {
			// body finished
			return nil
		}

		if token.TType == c.IF {
			cond := new(ifCondition)
			err := cond.parseSelf(data[i:])
			if err != nil {
				return err
			}
		}

		// check for function calls
		if token.TType == c.IDENT {
			// one char for the IDENT
			// one char for the RBRACE
			if i >= len(data)-2 {
				return errors.New("function body can contain a lone Identifier")
			}

			// check for method call
			// name(
			if data[i+1].TType == c.LPAREN {
				call := new(methodCall)
				err := call.parseSelf(data[i:])
				if err != nil {
					return err
				}

				b.Blocks = append(b.Blocks, call)
			}
		}
	}
	return errors.New("unterminated function body")
}

type ifCondition struct {
	expr   Statements
	ifBody []Statements
}

func (i *ifCondition) parseSelf(data []*lexer.Token) error {
	//TODO: implement me
	return nil
}

func (i *ifCondition) String() string {
	return fmt.Sprintf("ifCondition<%s>", i.expr)
}

type methodCall struct {
	name string
	args []Statements
}

func (m *methodCall) parseSelf(data []*lexer.Token) error {
	// name(args)
	m.args = make([]Statements, 0, 3)
	m.name = data[0].Lex
	data = data[2:]

	for i := 0; i < len(data); i++ {
		token := data[i]

		// successfully parsed all args
		if token.TType == c.RPAREN {
			return nil
		}

		// skip commas
		if token.TType == c.COMMA {
			continue
		}

		m.args = append(m.args,
			&refVar{
				name: token.Lex,
			})

	}
	return errors.New("unterminated function call")
}

func (m *methodCall) String() string {
	return fmt.Sprintf("%s(%s)", m.name, m.args)
}

// refVar
//
// reference variable for example in func calls
// add(refVar, refVar)
type refVar struct {
	name string
}

func (r *refVar) parseSelf(data []*lexer.Token) error {
	return nil
}

func (r *refVar) String() string {
	return fmt.Sprintf("refVar=%s", r.name)
}

// ArgVariable
//
// var that is in function definitions
// int main(ArgVariable, ArgVariable)
type ArgVariable struct {
	varType string
	name    string
}

// not necesarry it's easier to just set the values explicitly
func (v *ArgVariable) parseSelf(data []*lexer.Token) error {
	return nil
}

func (v *ArgVariable) String() string {
	return fmt.Sprintf("%s(%s)", v.name, v.varType)
}

type returnStmt struct {
	retExpr Statements
}

func (ret *returnStmt) parseSelf(data []*lexer.Token) error {
	//TODO implement me
	panic("implement me")
}

func (ret *returnStmt) String() string {
	return fmt.Sprintf("return %s", ret.retExpr)
}

type OpExpr struct {
	Left, Right Statements
	Op          c.TokenType
}

func (opE *OpExpr) parseSelf(data []*lexer.Token) error {
	//TODO implement me
	panic("implement me")
}

func (opE *OpExpr) String() string {
	return fmt.Sprintf("%s %s %s)", opE.Right, opE.Op, opE.Left.String())
}

type IntegerLit struct {
	Value int16
}

func (lit *IntegerLit) parseSelf(data []*lexer.Token) error {
	//TODO implement me
	panic("implement me")
}

func (lit *IntegerLit) String() string {
	return fmt.Sprintf("%d", lit.Value)
}

// type FloatLit struct{}

type StringLit struct {
	Value string
}

func (lit *StringLit) parseSelf(data []*lexer.Token) error {
	//TODO implement me
	panic("implement me")
}

func (lit *StringLit) String() string {
	return fmt.Sprintf("\"%s\"", lit.Value)
}

type Parser struct {
	AST *Statements
}
