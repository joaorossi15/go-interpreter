package parser

import (
	"fmt"
	"strconv"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_ int = iota // assign values 1 to 7 for the constants to get precedence
	LOWEST
	EQUALS      // ==
	LESSGREATER // > OR <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X OR !X
	CALL        // MyFunction(x)
)

type (
	// only right side of exp
	prefixParseFn func() ast.Expression

	// left side of expression is passed as input
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) (p *Parser) {
	p = &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.regPrefix(token.IDENT, p.parseIdentifier)
	p.regPrefix(token.INT, p.parseInteger)
	p.regPrefix(token.BANG, p.parsePrefixExpression)
	p.regPrefix(token.MINUS, p.parsePrefixExpression)

	p.nextToken() // initializes next token
	p.nextToken() // initializes curr token

	return
}

func (p *Parser) regPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) regInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tk token.TokenType) {
	err := fmt.Sprintf("expected next token to be %s, got %s instead", tk, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		st := p.parseStatement()

		if st != nil {
			program.Statements = append(program.Statements, st)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	st := &ast.LetStatement{Token: p.curToken}

	// verify if token type is IDENTIFIER
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// create identifier node
	st.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// TODO: parse expressions
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return st
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	st := &ast.ReturnStatement{Token: p.curToken}

	// TODO parse expressions
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}
	return st
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	st := &ast.ExpressionStatement{Token: p.curToken}
	st.Expression = p.parseExpression(LOWEST) // we pass the lowest precedence operator because we didnt parse anything yet, so we cant compare precedence

	// we dont use expectPeek() because for expressions the semicolon is optional
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return st
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	intLiteral := &ast.IntegerLiteral{Token: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as int", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	intLiteral.Value = val

	return intLiteral
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// creates a prefix operation node
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken() // advance the token pointer to get the integer or identifier

	// calls parseExpression to parse int or identifier and complete the node
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) noPrefixParseError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type] // returns function associated with token type

	if prefix == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}

	leftExpression := prefix() // calls function associated with token type

	return leftExpression
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}
