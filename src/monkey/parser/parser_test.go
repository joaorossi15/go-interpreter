package parser

import (
	"fmt"
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

func TestIntegerExpression(t *testing.T) {
	input := "5;"

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

	il, ok := st.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, instead=%T", st.Expression)
	}

	if il.Value != 5 {
		t.Errorf("ident.Value not 5, got=%v", il.Value)
	}

	if il.TokenLiteral() != "5" {
		t.Errorf("ident.Value not 5, got=%v", il.TokenLiteral())
	}
}

func TestPrefixOperator(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		par := NewParser(l)
		prog := par.ParseProgram()
		checkParserErrors(t, par)

		if len(prog.Statements) != 1 {
			t.Fatalf("program statements does not contain 1 statement")
		}

		st, ok := prog.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStratement, instead=%T", prog.Statements[0])
		}

		expression, ok := st.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral, instead=%T", st.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, tt.intValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, right ast.Expression, value int64) bool {
	integer, ok := right.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("right not *ast.IntegerLiteral. got=%T", right)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d, got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
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
