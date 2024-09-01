package ast

import (
	"bytes"

	"github.com/ollybritton/calclang/token"
)

// VariableAssignment represents the process of assignment to a variable in the AST.
// Example: `10 -> a`.
type VariableAssignment struct {
	Tok   token.Token // the token.ASSIGN token.
	Name  *Identifier
	Value Expression
}

func (va *VariableAssignment) statementNode()     {}
func (va *VariableAssignment) Token() token.Token { return va.Tok }
func (va *VariableAssignment) String() string {
	var out bytes.Buffer

	out.WriteString(va.Value.String())
	out.WriteString(" -> ")
	out.WriteString(va.Name.String())

	return out.String()
}

// InputAssignment represents statements like "? -> A" in the AST.
type InputAssignment struct {
	Tok  token.Token
	Name *Identifier
}

func (ia *InputAssignment) statementNode()     {}
func (ia *InputAssignment) Token() token.Token { return ia.Tok }
func (ia *InputAssignment) String() string {
	var out bytes.Buffer

	out.WriteString("? -> ")
	out.WriteString(ia.Name.String())

	return out.String()
}

// ExpressionStatement is a single expression by itself on one line.
// Example: `{start} a+10 {end}` (where start & end are the start and end of the line)
// General: `{start}{expression}{end}`
type ExpressionStatement struct {
	Tok        token.Token // The first token of the expression.
	Expression Expression
}

func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) Token() token.Token { return es.Tok }
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}
