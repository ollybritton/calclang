package parser

import "github.com/ollybritton/calclang/token"

// Definitions of operator precedences.
const (
	_ int = iota
	LOWEST
	SUM     // + or -
	PRODUCT // * or /
	PREFIX  // -X
	CALL    // fn(x)
)

// Mappings of precedences to their token types.
var precedences = map[token.Type]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}
