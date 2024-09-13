package object

import (
	"bytes"
	"fmt"
	"strings"

	"monkey/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	ERROR_OBJ        = "ERROR"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	NULL_OBJ         = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type String struct {
	Value string
}

type Return struct {
	Value Object
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Enviroment
}

type Error struct {
	Value string
}

type Null struct{}

type Enviroment struct {
	store map[string]Object
	outer *Enviroment
}

func NewEnclosedEnviroment(outer *Enviroment) *Enviroment {
	env := NewEnviroment()
	env.outer = outer
	return env
}

func NewEnviroment() *Enviroment {
	s := make(map[string]Object)
	return &Enviroment{store: s, outer: nil}
}

func (e *Enviroment) Value(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Value(name)
	}
	return obj, ok
}

func (e *Enviroment) Add(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

func (rt *Return) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rt *Return) Inspect() string  { return rt.Value.Inspect() }

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return "ERROR: " + err.Value }

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
