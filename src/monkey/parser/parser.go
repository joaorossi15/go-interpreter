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

// precedence table to map token type to precedence
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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
	p.regPrefix(token.TRUE, p.parseBoolean)
	p.regPrefix(token.FALSE, p.parseBoolean)
	p.regPrefix(token.LPAREN, p.parseGroupedExpressions)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.regInfix(token.PLUS, p.parseInfixExpression)
	p.regInfix(token.MINUS, p.parseInfixExpression)
	p.regInfix(token.SLASH, p.parseInfixExpression)
	p.regInfix(token.ASTERISK, p.parseInfixExpression)
	p.regInfix(token.EQ, p.parseInfixExpression)
	p.regInfix(token.NOT_EQ, p.parseInfixExpression)
	p.regInfix(token.LT, p.parseInfixExpression)
	p.regInfix(token.GT, p.parseInfixExpression)

	p.nextToken() // initializes next token
	p.nextToken() // initializes curr token

	return
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	st.Value = p.parseExpression(LOWEST)

	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return st
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	st := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	st.Value = p.parseExpression(LOWEST)

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

func (p *Parser) parseBoolean() ast.Expression {
	value := false
	if p.curToken.Type == token.TRUE {
		value = true
	}
	return &ast.Boolean{Token: p.curToken, Value: value}
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpressions() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		return nil
	}
	p.nextToken()
	return expression
}

func (p *Parser) noPrefixParseError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	/*
	   case 1 + 2 + 3
	       ast needs two infix expression nodes
	       first node:
	           - right node: 3
	           - operator: +
	           - left node: infix 1 + 2
	       second node:
	           - right node: 2
	           - operator: +
	           - left node: 1

	       tree:
	               ast.infixexpression
	               |                |
	       ast.infixexpression    ast.integer
	       |                 |         |
	   ast.integer     ast.integer     3
	       |                 |
	       1                 2

	   code working:
	       - parseExpression check if there is prefix function associated with curToken, and there is for 1 (INT);
	       - leftExp := *ast.IntegerLiteral
	       - for loop checks that peekToken is not semicolon and peekPrecedence is smaller
	       - inside for loop, fetch infixParse function assigned to next token (+);
	       - before executes, advances token so cur = + and peek = 2
	       - inside parse infix, creates ast.InfixExpression with operator = + and left = 1 (ast.IntegerLiteral) already defined
	       - saves the + precedence
	       - advances token
	       - calls parseExpression to next expression using the + precedence
	       - now parseExpression is called with cur = 2 and peek = +
	       - parse 2 as ast.IntegerLiteral
	       - for loop doesnt execute, because precedence of + (the argument passed as precedence) is not smalller than of peek (+), so 2 is returned
	       - it goes back to parseInfix and the right node receives 2, constructing the ast.InfixExpression for 1 + 2
	       - everything executes again, but left node is now ast.InfixExpression of 1 + 2 and right node is ast.IntegerLiteral of 3

	*/
	prefix := p.prefixParseFns[p.curToken.Type] // returns function associated with token type

	if prefix == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}

	leftExpression := prefix() // calls function associated with token type and returns prefix expression in the form of ast.PrefixExpression

	// check if next token is not semicolon and current operator has
	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		// get infix function
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			p.noPrefixParseError(p.curToken.Type)
			return nil
		}
		p.nextToken() // advances token so that we can parse the new
		leftExpression = infix(leftExpression)
	}

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
