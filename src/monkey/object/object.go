package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	ERROR_OBJ        = "ERROR"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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

type Return struct {
	Value Object
}

type Error struct {
	Value string
}

type Null struct{}

type Enviroment struct {
	store map[string]Object
}

func NewEnviroment() *Enviroment {
	s := make(map[string]Object)
	return &Enviroment{store: s}
}

func (e *Enviroment) Value(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Enviroment) Add(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

func (rt *Return) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rt *Return) Inspect() string  { return rt.Value.Inspect() }

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return "ERROR: " + err.Value }

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
