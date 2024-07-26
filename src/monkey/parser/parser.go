package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) (p *Parser) {
	p = &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken() // initializes next token
	p.nextToken() // initializes curr token

	return
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

func (par *Parser) PrintParser(prog *ast.Program) {
	for i, p := range prog.Statements {
		if p != nil {
			switch p.(type) {
			case (*ast.ReturnStatement):
				fmt.Printf("Node: %v, Node value: %q\n", i, p.TokenLiteral())
			default:
				if p.(*ast.LetStatement) == nil {
					fmt.Printf("Error: %q\n", par.Errors())
				} else {
					fmt.Printf("Node: %v, Node value: %q\n", i, p.TokenLiteral())
				}
			}
		}
	}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	st := &ast.LetStatement{Token: p.curToken}

	// verify if token type is IDENTIFIER
	if !p.expectPeek(token.IDENT) {
		return nil
	}

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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}
