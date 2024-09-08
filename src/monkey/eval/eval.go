package eval

import (
	"fmt"

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

func Eval(node ast.Node, env *object.Enviroment) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program:
		return evalProgram(node, env) // calls evalStatements for all of the statements

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.Return{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Add(node.Name.Value, val)

		// expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, node.Operator, right)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.FunctionLiteral:
		params := node.Arguments
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToObj(node.Value)
	}
	return nil
}

func evalProgram(p *ast.Program, env *object.Enviroment) object.Object {
	var result object.Object

	for _, st := range p.Statements {
		result = Eval(st, env)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Enviroment) object.Object {
	var result object.Object

	for _, st := range block.Statements {
		result = Eval(st, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Enviroment) object.Object {
	val, ok := env.Value(node.Value)

	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

func evalExpressions(expressions []ast.Expression, env *object.Enviroment) []object.Object {
	var res []object.Object

	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}
	return res
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendedFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapedReturnValue(evaluated)
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Enviroment {
	env := object.NewEnclosedEnviroment(fn.Env)
	for idx, param := range fn.Parameters {
		env.Add(param.Value, args[idx])
	}
	return env
}

func unwrapedReturnValue(obj object.Object) object.Object {
	if ret, ok := obj.(*object.Return); ok {
		return ret.Value
	}
	return obj
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
	} else if fn, ok := BOOLOPERATIONS[operator]; ok {
		return nativeBoolToObj(fn(valueLeft, valueRight))
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIfExpression(node *ast.IfExpression, env *object.Enviroment) object.Object {
	condition := Eval(node.Condition, env)

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func nativeBoolToObj(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Value: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
