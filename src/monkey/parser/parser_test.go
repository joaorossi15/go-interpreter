package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

func TestReturnStatement(t *testing.T) {
	input := `
    return 10;
    return 5;
    return 1201201;
    `

	l := lexer.New(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Program does not contain 3 statements")
	}

	for _, st := range program.Statements {
		returnSt, ok := st.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("s not *ast.LetStatement. got=%T", st)
		}
		if returnSt.TokenLiteral() != "return" {
			t.Errorf("token literal not let, got %q", st.TokenLiteral())
		}
	}
}

func TestLetStatement(t *testing.T) {
	input := `
    let x = 5;
    let y = 10;
    let foobar = 838383;
    `

	l := lexer.New(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program does not contain 3 statements")
	}

	tests := []struct {
		expectedId string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		st := program.Statements[i]
		if !testLetStatement(t, st, tt.expectedId) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, expected string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("token literal not let, got %q", s.TokenLiteral())
		return false
	}

	letSt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Fatalf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letSt.Name.Value != expected {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", expected, letSt.Name.Value)
		return false
	}

	if letSt.Name.TokenLiteral() != expected {
		t.Errorf("s.Name not '%s'. got=%s", expected, letSt.Name)
		return false
	}
	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	par := NewParser(l)
	prog := par.ParseProgram()
	checkParserErrors(t, par)

	if len(prog.Statements) != 1 {
		t.Fatalf("program has %d statements", len(prog.Statements))
	}

	st, ok := prog.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStratement, instead=%T", prog.Statements[0])
	}

	id, ok := st.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier, instead=%T", st.Expression)
	}

	if id.Value != "foobar" {
		t.Errorf("ident.Value not foobar, got=%s", id.Value)
	}

	if id.TokenLiteral() != "foobar" {
		t.Errorf("ident.Value not foobar, got=%s", id.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}

	t.FailNow()
}
