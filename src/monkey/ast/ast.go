package ast

import "monkey/token"

// simple node interface
type Node interface {
	TokenLiteral() string
}

// node that represents a statemet (let, return, etc)
type Statement interface {
	Node
	statementNode()
}

// node that represents an expression (value)
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // hold x in let x = 5;
	Value Expression  // expression that produces the value, 5 in let x = 5;
}

type ReturnStatement struct {
	Token token.Token // token.RETURN
	Value Expression  // expression that is returned
}

// hold x in let x = 5;
type Identifier struct {
	Token token.Token // name of the identifier, token.IDENT
	Value string      // value of the identifier, in this case x
}

func (lt *LetStatement) statementNode() {}

func (lt *LetStatement) TokenLiteral() string { return lt.Token.Literal }

func (rt *ReturnStatement) statementNode() {}

func (rt *ReturnStatement) TokenLiteral() string { return rt.Token.Literal }

func (id *Identifier) expressionNode() {}

func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

// returns the root of the program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
