package eval

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE       = &object.Boolean{Value: true}
	FALSE      = &object.Boolean{Value: false}
	NULL       = &object.Null{}
	OPERATIONS = map[string]func(int64, int64) int64{
		"+": func(a, b int64) int64 { return a + b },
		"-": func(a, b int64) int64 { return a - b },
		"*": func(a, b int64) int64 { return a * b },
		"/": func(a, b int64) int64 { return a / b },
	}
	BOOLOPERATIONS = map[string]func(int64, int64) bool{
		">":  func(a, b int64) bool { return a > b },
		"<":  func(a, b int64) bool { return a < b },
		"!=": func(a, b int64) bool { return a != b },
		"==": func(a, b int64) bool { return a == b },
	}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program:
		return evalStatements(node.Statements) // calls evalStatements for all of the statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

		// expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(left, node.Operator, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToObj(node.Value)
	}
	return nil
}

func evalStatements(sts []ast.Statement) object.Object {
	var result object.Object

	for _, st := range sts {
		result = Eval(st)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL
	}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)
	case operator == "==":
		return nativeBoolToObj(left == right)
	case operator == "!=":
		return nativeBoolToObj(left != right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	valueLeft := left.(*object.Integer).Value
	valueRight := right.(*object.Integer).Value
	if fn, ok := OPERATIONS[operator]; ok {
		return &object.Integer{Value: fn(valueLeft, valueRight)}
	}

	return nativeBoolToObj(BOOLOPERATIONS[operator](valueLeft, valueRight))
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func nativeBoolToObj(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}
