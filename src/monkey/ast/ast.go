package ast

import (
	"bytes"

	"monkey/token"
)

// simple node interface
type Node interface {
	TokenLiteral() string
	String() string // method for printing and debugging
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

// statement that consists of only a single expression and works like a wrapper
// we need this because monkey is a script language, so lines with only a expressions statement are legal
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// hold x in let x = 5;
type Identifier struct {
	Token token.Token // name of the identifier, token.IDENT
	Value string      // value of the identifier, in this case x
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (lt *LetStatement) statementNode()       {}
func (lt *LetStatement) TokenLiteral() string { return lt.Token.Literal }

func (lt *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(lt.TokenLiteral() + " ")
	out.WriteString(lt.Name.String())
	out.WriteString(" = ")

	if lt.Value != nil {
		out.WriteString(lt.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (rt *ReturnStatement) statementNode()       {}
func (rt *ReturnStatement) TokenLiteral() string { return rt.Token.Literal }

func (rt *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rt.TokenLiteral())

	if rt.Value != nil {
		out.WriteString(rt.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (id *Identifier) expressionNode()      {}
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

func (id *Identifier) String() string { return id.Value }

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// return whole expression as string
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// returns the root of the program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// program string method, writing the value of each satement String() method and returning it
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
