package parser

import (
	"strconv"

	"github.com/ollybritton/calclang/ast"
	"github.com/ollybritton/calclang/lexer"
	"github.com/ollybritton/calclang/token"
)

// Definition of parsing functions.
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser represents a parser for an calclang program.
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []error

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New returns a new parser from a given lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = map[token.Type]prefixParseFn{
		token.IDENT: p.parseIdentifier,
		token.INT:   p.parseIntegerLiteral,
		token.FLOAT: p.parseFloatLiteral,

		token.MINUS: p.parsePrefixExpression,

		token.LPAREN: p.parseGroupedExpression,
	}

	p.infixParseFns = map[token.Type]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
	}

	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns the errors that occured during parsing.
func (p *Parser) Errors() []error {
	return p.errors
}

// addError adds an error to the parser's internal error list.
func (p *Parser) addError(err error) {
	p.errors = append(p.errors, err)
}

// Parse parses the input program into a ast.Program.
func (p *Parser) Parse() *ast.Section {
	program := &ast.Section{}

	for {
		if p.curTokenIs(token.EOF) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	for p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	if p.curTokenIs(token.COLON) {
		p.nextToken()
	}

	if p.curTokenIs(token.EOF) {
		return nil
	}

	switch p.curToken.Type {
	case token.LBRACE: // Repeat
		return p.parseRepeatStatement()
	default:
		expr := p.parseExpression(LOWEST)
		var stmt ast.Statement

		if p.peekTokenIs(token.ASSIGN_TO) {
			p.nextToken() // current token is now ->
			p.nextToken() // current token is now the start of identifier

			ident, _ := p.parseIdentifier().(*ast.Identifier)
			stmt = &ast.VariableAssignment{
				Tok:   p.curToken,
				Name:  ident,
				Value: expr,
			}
		} else {
			stmt = &ast.ExpressionStatement{
				Tok:        expr.Token(),
				Expression: expr,
			}
		}

		return stmt
	}

}

// Expression Parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addError(
			NewNoPrefixParseFnError(p.curToken, p.peekToken, p.curToken.Type),
		)

		return nil
	}
	leftExp := prefix()

	for !(p.peekTokenIs(token.EOF)) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Tok: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(
			NewFloatParseError(p.curToken, p.peekToken, p.curToken.Literal),
		)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Tok: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(
			NewFloatParseError(p.curToken, p.peekToken, p.curToken.Literal),
		)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	// some expression ( 1, 2 )
	//                 ^

	return p.parseSubroutineCall(left)
}

func (p *Parser) parseSubroutineCall(expression ast.Expression) ast.Expression {
	exp := &ast.SubroutineCall{Tok: p.curToken, Subroutine: expression}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Tok:      p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	var expression = &ast.InfixExpression{
		Tok:      p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseRepeatStatement() ast.Statement {
	repeat := &ast.RepeatStatement{Tok: p.curToken}

	p.nextToken()
	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}

		if p.curTokenIs(token.RBRACE) {
			return repeat
		}

		stmt := p.parseStatement()
		if stmt != nil {
			repeat.Statements = append(repeat.Statements, stmt)
		}

		p.nextToken()
	}

	return repeat
}
