package lexer

import (
	"testing"

	"github.com/ollybritton/calclang/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `5 -> A:1.2 -> beans
tanh(10*pi)
:::
? -> ident`

	tests := []token.Token{
		{Type: token.INT, Literal: "5", Line: 0, StartCol: 0},
		{Type: token.ASSIGN_TO, Literal: "->", Line: 0, StartCol: 2},
		{Type: token.IDENT, Literal: "A", Line: 0, StartCol: 5},
		{Type: token.COLON, Literal: ":", Line: 0},
		{Type: token.FLOAT, Literal: "1.2", Line: 0},
		{Type: token.ASSIGN_TO, Literal: "->", Line: 0},
		{Type: token.IDENT, Literal: "beans", Line: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 0},
		{Type: token.IDENT, Literal: "tanh", Line: 1},
		{Type: token.LPAREN, Literal: "(", Line: 1},
		{Type: token.INT, Literal: "10", Line: 1},
		{Type: token.ASTERISK, Literal: "*", Line: 1},
		{Type: token.IDENT, Literal: "pi", Line: 1},
		{Type: token.RPAREN, Literal: ")", Line: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 1},
		{Type: token.TRIPLE_COLON, Literal: ":::", Line: 2},
		{Type: token.NEWLINE, Literal: "\n", Line: 2},
		{Type: token.QUESTION_MARK, Literal: "?", Line: 3},
		{Type: token.ASSIGN_TO, Literal: "->", Line: 3},
		{Type: token.IDENT, Literal: "ident", Line: 3},
		{Type: token.EOF, Literal: "", Line: 3},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.Type, tok.Type, "token type wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.Literal, tok.Literal, "token literal wrong for token %s, expecting %s", tok, tt.String())

		if tt.Line != 0 {
			assert.Equal(t, tt.Line, tok.Line, "token line number wrong for token %s, expecting %s", tok, tt.String())
		}

		if tt.StartCol != 0 {
			assert.Equal(t, tt.StartCol, tok.StartCol, "token StartCol number wrong for token %s, expecting %s", tok, tt.String())
		}

		if tt.EndCol != 0 {
			assert.Equal(t, tt.EndCol, tok.EndCol, "token EndCol number wrong for token %s, expecting %s", tok, tt.String())
		}
	}

	assert.Equal(t, byte(0), l.peekChar(), "lexer should have read all input before tests finish, not enough test cases")
}
