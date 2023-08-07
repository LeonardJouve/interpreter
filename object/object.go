package object

import "fmt"

type ObjectType string

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

type Null struct{}

type Return struct {
	Value Object
}

type Error struct {
	Value string
}

const (
	NULL    = "NULL"
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	RETURN  = "RETURN"
	ERROR   = "ERROR"
)

func (integer *Integer) Type() ObjectType {
	return INTEGER
}
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

func (boolean *Boolean) Type() ObjectType {
	return BOOLEAN
}
func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

func (null *Null) Type() ObjectType {
	return NULL
}
func (null *Null) Inspect() string {
	return "null"
}

func (ret *Return) Type() ObjectType {
	return RETURN
}
func (ret *Return) Inspect() string {
	return ret.Value.Inspect()
}

func (err *Error) Type() ObjectType {
	return ERROR
}
func (err *Error) Inspect() string {
	return "[Error] " + err.Value
}
