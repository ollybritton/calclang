package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ollybritton/calclang/token"
)

// Identifier represents an identifier in the AST. Idents are expressions because they
// produce values (the value they represent)
type Identifier struct {
	Tok      token.Token // the token.IDENT token.
	Constant bool
	Value    string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) Token() token.Token { return i.Tok }
func (i *Identifier) String() string {
	var out bytes.Buffer

	out.WriteString(i.Value)

	return out.String()
}

// IntegerLiteral represents an integer value in the AST.
// Example: `5`
// General: `{token.INT}`
type IntegerLiteral struct {
	Tok   token.Token // the token.INT token.
	Value int64
}

func (il *IntegerLiteral) expressionNode()    {}
func (il *IntegerLiteral) Token() token.Token { return il.Tok }
func (il *IntegerLiteral) String() string {
	return fmt.Sprint(il.Value)
}

// FloatLiteral represents an float value in the AST.
// Example: `5.5`
// General: `{token.INT}`
type FloatLiteral struct {
	Tok   token.Token // the token.FLOAT token.
	Value float64
}

func (fl *FloatLiteral) expressionNode()    {}
func (fl *FloatLiteral) Token() token.Token { return fl.Tok }
func (fl *FloatLiteral) String() string {
	return fmt.Sprint(fl.Value)
}

// PrefixExpression represents an expression involving a prefix operator.
// Example: `-10`
// General: `{- or !}{expression}`
type PrefixExpression struct {
	Tok      token.Token // the token of the prefix operator.
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()    {}
func (pe *PrefixExpression) Token() token.Token { return pe.Tok }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an expression involving an operator sandwiched between
// two other expressions.
// Example: `10-5`
// General: `{expression}{opeator}{expression}`
type InfixExpression struct {
	Tok token.Token

	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()    {}
func (ie *InfixExpression) Token() token.Token { return ie.Tok }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// SubroutineCall represents a call to a subroutine within the AST.
// Example: `add(1,2)`
// General: `{IDENT}({expression}, {expression}...)`
type SubroutineCall struct {
	Tok        token.Token // The '(' token
	Subroutine Expression
	Arguments  []Expression
}

func (sc *SubroutineCall) expressionNode()    {}
func (sc *SubroutineCall) Token() token.Token { return sc.Tok }
func (sc *SubroutineCall) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range sc.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(sc.Subroutine.String())
	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}
